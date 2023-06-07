-- name: CreateUserToken :one
INSERT INTO user_tokens (owner, token, expire_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserToken :one
SELECT *
FROM user_tokens
WHERE owner = $1 AND id = $2;

-- name: GetUserTokens :many
SELECT *
FROM user_tokens
WHERE owner = $1;

-- name: DeleteUserToken :exec
DELETE FROM user_tokens
WHERE id = $1 AND owner = $2;