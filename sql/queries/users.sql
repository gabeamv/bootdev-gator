-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE $1 = users.name;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
