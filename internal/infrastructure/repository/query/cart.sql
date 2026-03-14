-- name: SaveCart :exec
INSERT INTO carts (id, buyer_id, created_at) VALUES ($1, $2, $3);

-- name: FindCartByBuyerID :one
SELECT * FROM carts WHERE buyer_id = $1;

-- name: SaveCartItem :exec
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES ($1, $2, $3)
    ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = EXCLUDED.quantity;

-- name: DeleteCartItem :exec
DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2;

-- name: FindCartItems :many
SELECT * FROM cart_items WHERE cart_id = $1;
