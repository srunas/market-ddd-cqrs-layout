-- name: SaveUser :exec
INSERT INTO users (id, username, surname, role, email, created_at, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, surname = $3, email = $4, enabled = $5
WHERE id = $1;
