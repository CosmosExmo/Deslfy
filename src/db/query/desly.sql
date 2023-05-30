-- name: CreateDesly :one
INSERT INTO deslies (redirect, desly)
VALUES ($1, substr(gen_random_uuid()::text, 1, 6))
RETURNING *;

/* -- name: GetDesly :one
SELECT *
FROM deslies
WHERE id = $1
LIMIT 1; */

-- name: GetDeslyByDesly :one
SELECT *
FROM deslies
WHERE desly = $1
LIMIT 1;