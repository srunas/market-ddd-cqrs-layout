-- name: SaveOrder :exec
INSERT INTO orders (id, buyer_id, status, total, currency, payment_method, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: FindOrdersByBuyerID :many
SELECT * FROM orders WHERE buyer_id = $1;

-- name: UpdateOrder :exec
UPDATE orders SET status = $2 WHERE id = $1;

-- name: SaveOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity, price_at_order)
VALUES ($1, $2, $3, $4);

-- name: FindOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;
