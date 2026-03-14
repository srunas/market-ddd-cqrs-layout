-- name: SaveCategory :exec
INSERT INTO categories (id, name, parent_id, level, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: FindCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: FindCategoryByIDForUpdate :one
SELECT * FROM categories WHERE id = $1 FOR UPDATE;

-- name: FindAllCategories :many
SELECT * FROM categories;
