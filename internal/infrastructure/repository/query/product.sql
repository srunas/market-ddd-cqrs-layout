-- name: SaveProduct :exec
INSERT INTO products (id, name, description, price, currency, stock, seller_id, created_at, updated_at, active, attributes, category_ids)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: FindProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: UpdateProduct :exec
UPDATE products
SET name = $2, description = $3, price = $4, stock = $5, updated_at = $6, active = $7, category_ids = $8
WHERE id = $1;

-- name: GetProductList :many
SELECT * FROM products
WHERE
    (sqlc.narg('category_id')::uuid IS NULL OR sqlc.narg('category_id')::uuid = ANY(category_ids))
  AND (sqlc.narg('min_price')::numeric IS NULL OR price >= sqlc.narg('min_price')::numeric)
  AND (sqlc.narg('max_price')::numeric IS NULL OR price <= sqlc.narg('max_price')::numeric)
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
