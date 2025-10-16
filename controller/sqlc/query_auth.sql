-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (id, hashed_password)
VALUES (?, ?)
RETURNING *;

-- name: SetUserPassword :one
UPDATE users
set hashed_password = ?
WHERE id = ?
RETURNING *;

-- name: DeleteAuthor :exec
DELETE
FROM users
WHERE id = ?;