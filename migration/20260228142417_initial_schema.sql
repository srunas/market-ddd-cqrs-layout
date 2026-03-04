-- migrations/20250228123456_initial_schema.sql

-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Auth
CREATE TABLE auth (
                      id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      username   TEXT NOT NULL UNIQUE,
                      password   TEXT NOT NULL,
                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                      auth_at    TIMESTAMPTZ
);

-- Users
CREATE TYPE user_role AS ENUM ('ADMIN', 'BUYER');

CREATE TABLE users (
                       id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username   TEXT NOT NULL UNIQUE,
                       surname    TEXT,
                       role       user_role NOT NULL,
                       auth_at    TIMESTAMPTZ,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       email      TEXT NOT NULL UNIQUE,
                       enabled    BOOLEAN NOT NULL DEFAULT TRUE
);

-- Categories
CREATE TABLE categories (
                            id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            name            TEXT NOT NULL,
                            parent_id       UUID REFERENCES categories(id) ON DELETE SET NULL,
                            level           INTEGER NOT NULL DEFAULT 0,
                            status          BOOLEAN NOT NULL DEFAULT TRUE,
                            attribute_set_id INTEGER,
                            created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Products
CREATE TABLE products (
                          id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          name        TEXT NOT NULL,
                          description TEXT,
                          price       NUMERIC(19,4) NOT NULL,
                          created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          currency    TEXT NOT NULL CHECK (currency IN ('USD', 'EUR')),
                          stock       INTEGER NOT NULL DEFAULT 0,
                          seller_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          active      BOOLEAN NOT NULL DEFAULT TRUE,
                          attributes  JSONB DEFAULT '{}',
                          category_ids UUID[] DEFAULT '{}'
);

-- Carts + items
CREATE TABLE carts (
                       id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       buyer_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cart_items (
                            cart_id    UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
                            product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                            quantity   INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
                            PRIMARY KEY (cart_id, product_id)
);

-- Orders + items
CREATE TYPE order_status AS ENUM ('CREATED', 'PROCESS', 'SUCCESS', 'FAIL', 'CANCELLED');
CREATE TYPE payment_method AS ENUM ('CARD', 'CASH', 'WALLET');

CREATE TABLE orders (
                        id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        buyer_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        status         order_status NOT NULL DEFAULT 'CREATED',
                        total          NUMERIC(19,4) NOT NULL DEFAULT 0,
                        currency       TEXT NOT NULL CHECK (currency IN ('USD', 'EUR')),
                        payment_method payment_method NOT NULL,
                        created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
                             order_id   UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                             quantity   INTEGER NOT NULL CHECK (quantity > 0),
                             price_at_order NUMERIC(19,4) NOT NULL,
                             PRIMARY KEY (order_id, product_id)
);

-- +goose Down

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS payment_method;

DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;

DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;

DROP TABLE IF EXISTS auth;

DROP EXTENSION IF EXISTS "uuid-ossp";
