# Marketplace — Go + DDD + CQRS

Учебный проект маркетплейса на Go. Цель — показать, как строится production-ready сервис с применением принципов **Domain-Driven Design (DDD)** и паттерна **CQRS**, с полноценным стеком наблюдаемости (метрики, логи, трейсы).

---

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Структура проекта](#структура-проекта)
3. [Архитектура и слои](#архитектура-и-слои)
4. [Почему DDD + CQRS](#почему-ddd--cqrs)
5. [Observability стек](#observability-стек)
6. [Инструменты](#инструменты)
7. [Запуск](#запуск)
8. [Нагрузочное тестирование](#нагрузочное-тестирование)
9. [Переменные окружения](#переменные-окружения)
10. [Используемые паттерны](#используемые-паттерны)
11. [Полезные ссылки](#полезные-ссылки)

---

## Обзор проекта

Маркетплейс — это платформа, где продавцы размещают товары, а покупатели их приобретают. Проект охватывает четыре основные доменные области:

| Домен | Описание |
|---|---|
| **Identity** | Регистрация, аутентификация, профиль пользователя |
| **Catalog** | Категории товаров, карточки товаров, публикация |
| **Cart** | Корзина покупателя: добавление, удаление, изменение количества |
| **Order** | Оформление заказа, статусная машина, отмена |

Технический стек: Go 1.25, PostgreSQL 15, pgx, sqlc, goose, JWT, OpenTelemetry, Grafana, Prometheus, Loki, Tempo.

---

## Структура проекта

```
market-ddd-cqrs-layout/
├── cmd/
│   └── app/
│       └── main.go                  # Точка входа, DI, запуск сервера
├── internal/
│   ├── domain/                      # Чистая бизнес-логика (ядро)
│   │   ├── entity/                  # Доменные сущности (Aggregate Roots)
│   │   │   ├── auth/
│   │   │   ├── cart/
│   │   │   ├── category/
│   │   │   ├── order/
│   │   │   ├── product/
│   │   │   └── user/
│   │   ├── repository/              # Интерфейсы репозиториев
│   │   ├── service/                 # Интерфейсы доменных сервисов
│   │   └── types/                   # Типизированные идентификаторы
│   ├── application/
│   │   └── service/                 # Реализации сервисов (use cases)
│   │       ├── identity-service/
│   │       ├── catalog-service/
│   │       ├── cart-service/
│   │       └── order-service/
│   ├── infrastructure/
│   │   └── repository/              # Реализации репозиториев (PostgreSQL)
│   │       ├── sqlcgen/             # Сгенерированный код (не трогать руками)
│   │       └── query/               # SQL-запросы для sqlc
│   ├── handler/                     # HTTP-обработчики
│   │   └── middleware/              # Middleware (логгер, метрики, трейсинг)
│   └── migrator/                    # Обёртка над goose
├── migration/                       # SQL-миграции
├── docker/                          # Конфиги инфраструктуры
│   ├── grafana/
│   ├── loki/
│   ├── prometheus/
│   └── tempo/
├── docker-compose.yml
├── sqlc.yaml
└── go.mod
```

---

## Архитектура и слои

Проект следует принципу **Clean Architecture**: зависимости направлены строго внутрь. Домен ничего не знает о базе данных, HTTP или внешних сервисах.

```
┌─────────────────────────────────────────┐
│             infrastructure              │  PostgreSQL, sqlc, внешние API
├─────────────────────────────────────────┤
│              application                │  Use cases, команды и запросы
├─────────────────────────────────────────┤
│                domain                   │  Бизнес-логика, сущности, правила
└─────────────────────────────────────────┘
             зависимости идут вверх →
```

### `internal/domain` — чистая бизнес-логика

Самый важный слой. Здесь живут правила предметной области. **Никаких SQL, HTTP, внешних библиотек** (кроме базовых вроде `errors`, `time`).

**Entity (Aggregate Roots)** — объекты с идентичностью и поведением:

- `entity/user` — пользователь (Username, Email, Role, Enabled)
- `entity/auth` — учётные данные (хэш пароля, JWT)
- `entity/product` — товар (Name, Price, Status, CategoryID)
- `entity/category` — категория с поддержкой вложенности (ParentID)
- `entity/cart` — корзина покупателя (Items, AddItem, RemoveItem)
- `entity/order` — заказ со статусной машиной (Created → Processed → Completed/Failed/Cancelled)

Пример доменной логики в `Order` — метод `Cancel` защищает инвариант: нельзя отменить завершённый или упавший заказ:

```go
func (o *Order) Cancel() error {
    if o.Status == StatusCompleted || o.Status == StatusFailed {
        return errors.New("cannot cancel completed or failed order")
    }
    o.Status = StatusCancelled
    return nil
}
```

**Repository interfaces** — контракты для работы с хранилищем. Домен описывает, что нужно, не зная как это реализовано:

```go
// internal/domain/repository/user.go
type UserRepository interface {
    Save(ctx context.Context, u *user.User) error
    FindByID(ctx context.Context, id types.UserID) (*user.User, error)
    FindByUsername(ctx context.Context, username string) (*user.User, error)
}
```

**Service interfaces** — контракты прикладных сервисов, которые вызывает HTTP-слой.

**Types** — типизированные идентификаторы (подробнее в разделе [Паттерны](#используемые-паттерны)).

### `internal/application` — команды и запросы (CQRS)

Здесь реализована бизнес-логика use cases: координация между репозиториями, транзакции, трейсинг. Каждый метод вынесен в **отдельный файл** — это упрощает навигацию и код-ревью.

```
application/service/
├── identity-service/
│   ├── implementation.go       # Структура сервиса, конструктор
│   ├── register_user.go        # Команда: зарегистрировать пользователя
│   ├── login_user.go           # Команда: войти в систему
│   └── get_user_profile.go     # Запрос: получить профиль
├── catalog-service/
│   ├── implementation.go
│   ├── create_product.go       # Команда: создать товар
│   ├── publish_product.go      # Команда: опубликовать товар
│   ├── create_category.go      # Команда: создать категорию
│   ├── list_products.go        # Запрос: список товаров
│   └── get_category_tree.go    # Запрос: дерево категорий
├── cart-service/
│   ├── implementation.go
│   ├── add_item.go             # Команда: добавить товар в корзину
│   ├── remove_item.go          # Команда: убрать товар из корзины
│   ├── decrease_quantity.go    # Команда: уменьшить количество
│   └── get_cart.go             # Запрос: получить корзину
└── order-service/
    ├── implementation.go
    ├── place_order.go          # Команда: оформить заказ из корзины
    ├── cancel_order.go         # Команда: отменить заказ
    ├── get_order.go            # Запрос: получить заказ по ID
    └── list_orders.go          # Запрос: список заказов пользователя
```

Каждый метод начинается с создания трейс-спана:

```go
func (s *Implementation) PlaceOrder(ctx context.Context, req service.PlaceOrderRequest) (
    service.PlaceOrderResponse, error) {
    ctx, span := prospan.Start(ctx)
    defer span.End()
    // ...бизнес-логика...
}
```

### `internal/infrastructure` — реализация репозиториев

Реализует интерфейсы из домена через PostgreSQL. Использует `sqlc` для type-safe работы с БД.

- `repository/*.go` — реализации репозиториев. Каждый принимает `*pgxpool.Pool` и оборачивает `sqlcgen.Queries`.
- `repository/sqlcgen/` — **сгенерированный код**, не редактировать вручную. Генерируется командой `sqlc generate`.
- `repository/query/` — SQL-запросы, из которых sqlc генерирует Go-код.

Каждый метод репозитория оборачивается в трейс-спан для видимости SQL-запросов в Tempo:

```go
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    ctx, span := prospan.Start(ctx)
    defer span.End()
    return r.q.SaveUser(ctx, ...)
}
```

### `internal/handler` — HTTP-обработчики

HTTP-хендлеры принимают интерфейсы сервисов, парсят запрос, вызывают сервис, возвращают JSON. Не содержат бизнес-логики.

**Цепочка middleware** (порядок важен):

```
otelhttp → WithLogger → WithRequestLogging → WithMetrics → mux
```

- `otelhttp` — создаёт корневой трейс-спан для каждого HTTP-запроса
- `WithLogger` — кладёт логгер в контекст
- `WithRequestLogging` — логирует метод, путь, статус, latency и `trace_id`
- `WithMetrics` — собирает Prometheus метрики (счётчики, гистограммы)

### `cmd/app/main.go` — точка входа

Выполняет роль **Composition Root** — единственное место, где собирается всё приложение:

1. Подключение к PostgreSQL через `pgxpool`
2. Инициализация JSON-логгера (zap)
3. Запуск провайдера OpenTelemetry трейсов → регистрация глобального провайдера
4. Применение миграций через goose
5. Создание transaction manager
6. Создание репозиториев
7. Создание сервисов с инъекцией репозиториев
8. Регистрация HTTP-маршрутов + `/metrics` endpoint
9. Запуск сервера с graceful shutdown (ловит `SIGTERM` / `SIGINT`)

---

## Почему DDD + CQRS

### DDD (Domain-Driven Design) простыми словами

Представьте, что вы строите маркетплейс. Без DDD вы скорее всего напишете функции типа `UpdateOrderStatus(id, status string)` — простая операция над базой данных. Проблема в том, что такой код позволяет перевести заказ из `completed` обратно в `created`, что в реальности невозможно.

DDD говорит: **бизнес-правила должны жить в коде**, а не только в головах разработчиков. Поэтому `Order` — это не просто структура с полями, а объект с поведением:

```go
// Нельзя просто поставить Status = "cancelled"
// Нужно вызвать метод, который проверит инварианты
err := order.Cancel() // вернёт ошибку если заказ уже завершён
```

### CQRS (Command Query Responsibility Segregation) простыми словами

CQRS — это разделение операций на две категории:

- **Команды** (Commands) — изменяют состояние. Примеры: `PlaceOrder`, `RegisterUser`, `AddItemToCart`. Возвращают минимум данных или ничего.
- **Запросы** (Queries) — читают данные и не меняют состояние. Примеры: `GetOrder`, `ListProducts`, `GetCart`.

Зачем это нужно? В реальном маркетплейсе чтений в разы больше, чем записей. CQRS позволяет оптимизировать их независимо: например, направлять запросы на read-реплику БД, кэшировать результаты, строить денормализованные представления.

В нашем проекте CQRS выражается на уровне именования: каждый файл в `application/service/` — это либо команда (`place_order.go`), либо запрос (`get_order.go`).

---

## Observability стек

Проект поставляется с полным стеком наблюдаемости. Три типа сигналов связаны через `trace_id`:

```
Приложение
    │
    ├── /metrics ──────────→ Prometheus (сбор метрик каждые 15с)
    ├── stdout (JSON) ──────→ Loki (сбор логов через Docker logging driver)
    └── OTLP gRPC :4317 ───→ Tempo (сбор трейсов)
                                    │
                             Grafana (единый UI)
                             localhost:3000
```

### Как три сигнала связаны между собой

```
1. Prometheus показывает всплеск ошибок (EPS вырос)
2. Grafana → Loki: ищем логи за это время → находим trace_id
3. Grafana → Tempo: вставляем trace_id → видим дерево вызовов
```

Каждый лог содержит `trace_id`, что позволяет переходить между сигналами.

### Метрики (Prometheus)

Приложение экспортирует два вида метрик на `/metrics`:

**Go runtime метрики** — собираются автоматически через `promhttp`:
- `go_goroutines` — количество горутин
- `go_memstats_alloc_bytes` — выделенная память
- `go_gc_duration_seconds` — паузы GC

**HTTP метрики** — собираются middleware `WithMetrics`:
- `http_requests_total{method, path, status}` — счётчик запросов
- `http_request_duration_seconds{method, path}` — гистограмма latency

Полезные PromQL запросы для Grafana:
```promql
# RPS по endpoint-ам
sum by (method, path) (rate(http_requests_total[1m]))

# EPS (errors per second) с разбивкой по ручкам
sum by (path, status) (rate(http_requests_total{status=~"[45].."}[1m]))

# Latency p95
histogram_quantile(0.95, sum by (le, path) (rate(http_request_duration_seconds_bucket[1m])))

# Топ ошибок за 5 минут
sort_desc(sum by (path, status) (increase(http_requests_total{status=~"[45].."}[5m])))
```

### Логи (Loki)

Приложение пишет структурированные JSON-логи через `zap`. Каждая строка содержит `trace_id`:

```json
{
  "level": "info",
  "msg": "входящий запрос",
  "method": "POST",
  "path": "/api/v1/register",
  "status": 201,
  "duration_ms": 111,
  "trace_id": "527f0b9fcf036def01bfaf0b9c823dd9"
}
```

Поиск в Grafana через LogQL:
```logql
# Все логи сервиса
{service="marketplace"} | json

# Только ошибки
{service="marketplace"} | json | level="error"

# Конкретный endpoint
{service="marketplace"} |= "/api/v1/register"

# По trace_id
{service="marketplace"} | json | trace_id="527f0b9fcf036def01bfaf0b9c823dd9"
```

### Трейсы (Tempo)

Каждый HTTP-запрос порождает дерево спанов. Спаны добавлены на двух уровнях:

**Application layer** — каждый use case:
```
RegisterUser [~100ms]
PlaceOrder [~150ms]
```

**Infrastructure layer** — каждый вызов к БД:
```
RegisterUser [~100ms]
  ├── AuthRepository.Save [~40ms]   ← INSERT в auth
  └── UserRepository.Save [~55ms]   ← INSERT в users
```

Поиск трейса в Grafana: Explore → Tempo → Search → Service Name: `marketplace`.

### Дашборд Grafana

Рекомендуемые панели для дашборда маркетплейса:

| Панель | Запрос | Тип |
|--------|--------|-----|
| RPS по endpoint-ам | `sum by (method, path) (rate(http_requests_total[1m]))` | Time series |
| EPS по endpoint-ам | `sum by (path, status) (rate(http_requests_total{status=~"[45].."}[1m]))` | Time series |
| Latency p95 | `histogram_quantile(0.95, ...)` | Time series |
| Goroutines | `go_goroutines` | Time series |
| Heap память | `go_memstats_alloc_bytes` | Time series |
| Логи | `{service="marketplace"} \| json` | Logs |

---

## Инструменты

### sqlc — генерация кода из SQL

`sqlc` читает ваши SQL-запросы и генерирует типобезопасный Go-код. Вместо ручного `rows.Scan(...)` вы получаете готовые методы.

**Как использовать:**

1. Пишете SQL-запрос в `internal/infrastructure/repository/query/*.sql`
2. Запускаете `sqlc generate`
3. Получаете готовый код в `internal/infrastructure/repository/sqlcgen/`

> Папку `sqlcgen/` не редактировать вручную — изменения будут перезаписаны при следующей генерации.

### goose — миграции базы данных

`goose` управляет версиями схемы БД. Файлы миграций находятся в папке `migration/` и именуются по шаблону `YYYYMMDDHHMMSS_название.sql`.

Миграции применяются **автоматически при старте** приложения через `internal/migrator/migrator.go`.

Для ручного управления:
```bash
goose -dir migration postgres "$DB_URI" up    # применить все миграции
goose -dir migration postgres "$DB_URI" down  # откатить последнюю
goose -dir migration postgres "$DB_URI" status
```

### go-transaction-manager — управление транзакциями

[avito-tech/go-transaction-manager](https://github.com/avito-tech/go-transaction-manager) позволяет выполнять несколько операций с БД атомарно, передавая транзакцию через контекст:

```go
err = s.txManager.Do(ctx, func(ctx context.Context) error {
    if err := s.userRepo.Save(ctx, userEntity); err != nil {
        return err
    }
    return s.authRepo.Save(ctx, authEntity)
    // если здесь ошибка — оба Save откатятся
})
```

### prospan (observer) — OpenTelemetry трейсинг

`github.com/not-for-prod/observer` — обёртка над OpenTelemetry SDK:

- `tracer.NewProvider(...)` — инициализация провайдера, регистрация глобального `otel.SetTracerProvider`
- `prospan.Start(ctx)` — создание дочернего спана внутри метода
- `logger/zap` — JSON-логгер

### otelhttp — автоматический трейсинг HTTP

`go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` создаёт корневой спан для каждого HTTP-запроса и добавляет его в контекст. Это позволяет `WithRequestLogging` читать `trace_id` и записывать его в лог.

---

## Запуск

### Предварительные требования

- Docker и Docker Compose
- Go 1.25+
- `sqlc`: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- `goose`: `go install github.com/pressly/goose/v3/cmd/goose@latest`

### Шаг 1: Создать файл окружения

Создайте `.env` в корне проекта:

```env
DB_URI=postgres://demo:demo@mypostgres:5432/marketplace?sslmode=disable
POSTGRES_USER=demo
POSTGRES_PASSWORD=demo
POSTGRES_DB=postgres
POSTGRES_PORT=5432
APP_PORT=8080
JWT_SECRET=your-secret-key-change-in-production
MIGRATIONS_DIR=migration
OTEL_EXPORTER_OTLP_ENDPOINT=tempo:4317
```

### Шаг 2: Создать БД маркетплейса

```bash
docker compose up postgres -d
docker exec -it mypostgres psql -U demo -c "CREATE DATABASE marketplace;"
```

### Шаг 3: Запустить все сервисы

```bash
docker compose up --build -d
```

Поднимет: PostgreSQL, приложение Go, Tempo, Loki, Prometheus, Grafana.

### Шаг 4: Открыть Grafana

[http://localhost:3000](http://localhost:3000) — логин `admin`, пароль из `.env`.

### Локальный запуск без Docker

```bash
docker compose up postgres tempo loki prometheus grafana -d

DB_URI=postgres://demo:demo@localhost:5432/marketplace?sslmode=disable \
MIGRATIONS_DIR=migration \
JWT_SECRET=dev-secret \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
go run ./cmd/app/main.go
```

### Health checks

```
GET /healthz/live   # liveness: 200 если процесс жив
GET /healthz/ready  # readiness: 200 если БД доступна
GET /metrics        # Prometheus метрики
```

### API endpoints

```
POST /api/v1/register   # регистрация: {"username","email","password","surname"}
POST /api/v1/login      # вход: {"username","password"} → {"token":"..."}
```

---

## Нагрузочное тестирование

Для нагрузки используется Apache Bench (`ab`):

```bash
# Установка (macOS)
brew install httpd

# Создать файл с телом запроса
echo '{"username":"loadtest","password":"secret123"}' > /tmp/login.json

# 1000 запросов, 50 одновременных
ab -n 1000 -c 50 -T "application/json" \
  -p /tmp/login.json \
  http://localhost:8080/api/v1/login
```

Во время теста в Grafana видно:
- **RPS** растёт до 100+ req/s
- **Goroutines** прыгают с ~23 до 400+
- **Heap память** растёт
- **Latency p95** показывает реальную нагрузочную задержку

---

## Переменные окружения

| Переменная | Описание | Пример |
|---|---|---|
| `DB_URI` | Connection string PostgreSQL | `postgres://user:pass@host:5432/db?sslmode=disable` |
| `JWT_SECRET` | Секрет для подписи JWT токенов | `your-secret-key` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Адрес Tempo для отправки трейсов | `tempo:4317` |
| `MIGRATIONS_DIR` | Путь до папки с миграциями | `migration` |
| `POSTGRES_USER` | Пользователь PostgreSQL | `demo` |
| `POSTGRES_PASSWORD` | Пароль PostgreSQL | `demo` |
| `POSTGRES_DB` | Имя БД по умолчанию | `postgres` |
| `POSTGRES_PORT` | Порт PostgreSQL на хосте | `5432` |
| `APP_PORT` | Порт приложения на хосте | `8080` |

---

## Используемые паттерны

### Typed IDs — защита от перепутанных аргументов

```go
type UserID    uuid.UUID
type OrderID   uuid.UUID
type ProductID uuid.UUID

// Ошибка компиляции: cannot use orderID (type types.OrderID) as types.UserID
func GetUserOrders(userID types.UserID) { ... }
```

### Repository pattern

Интерфейс репозитория в домене, реализация в infrastructure. Позволяет тестировать бизнес-логику с mock-репозиторием и легко менять хранилище.

### Один метод — один файл

Каждый use case в отдельном файле: `place_order.go`, `cancel_order.go`. Проще навигация, меньше конфликтов в git.

### Middleware chain

```go
h := otelhttp.NewHandler(
    middleware.WithLogger(logger)(
        middleware.WithRequestLogging()(
            middleware.WithMetrics()(mux),
        ),
    ),
    "http",
)
```

Каждый middleware решает одну задачу. Порядок важен: `otelhttp` создаёт спан первым, затем остальные могут читать `trace_id` из контекста.

### NewFromDB конструктор

Отделяет создание новой сущности (генерация ID, начальный статус) от восстановления из БД:

```go
order := order.New(buyerID, currency, method)  // новая
order := toOrderDomain(dbRow)                  // из БД
```

---

## Полезные ссылки

- **Grafana provisioning**: https://grafana.com/docs/grafana/latest/administration/provisioning/
- **Loki конфигурация**: https://grafana.com/docs/loki/latest/configuration/
- **Tempo конфигурация**: https://grafana.com/docs/tempo/latest/configuration/
- **Prometheus конфигурация**: https://prometheus.io/docs/prometheus/latest/configuration/configuration/
- **OpenTelemetry Go**: https://opentelemetry.io/docs/languages/go/
- [avito-tech/go-transaction-manager](https://github.com/avito-tech/go-transaction-manager)
- [sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc)
- [pressly/goose](https://github.com/pressly/goose)
- [shopspring/decimal](https://github.com/shopspring/decimal)
- [jackc/pgx](https://github.com/jackc/pgx)
