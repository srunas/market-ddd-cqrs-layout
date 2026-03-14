-- name: Save :exec
INSERT INTO auth (id, username, password, created_at, auth_at)
VALUES ($1, $2, $3, $4, $5);

-- name: FindByUsername :one
SELECT * FROM auth WHERE username = $1;

-- name: UpdateAuth :exec
UPDATE auth
SET auth_at = $2
WHERE username = $1;
