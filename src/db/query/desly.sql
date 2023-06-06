-- name: CreateDesly :one
INSERT INTO deslies (redirect, desly, owner)
VALUES ($1, substr(gen_random_uuid()::text, 1, 6), $2)
RETURNING *;

/* -- name: GetDesly :one
SELECT *
FROM deslies
WHERE id = $1
LIMIT 1; */

-- name: GetDesly :one
SELECT *
FROM deslies
WHERE desly = $1 AND owner = $2
LIMIT 1;

-- name: GetRedirectByDesly :one
SELECT redirect
FROM deslies
WHERE desly = $1
LIMIT 1;