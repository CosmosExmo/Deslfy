-- name: CreateDesly :one
INSERT INTO deslies (redirect, desly)
VALUES ($1, $2)
RETURNING *;

-- name: GetDesly :one
SELECT *
FROM deslies
WHERE id = $1
LIMIT 1;

-- name: GetRedirectByDesly :one
SELECT *
FROM deslies
WHERE desly = $1
LIMIT 1;