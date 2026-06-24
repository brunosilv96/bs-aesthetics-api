-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (customer_id, token_hash, expires_at, created_at)
VALUES ($1, $2, $3, NOW()) RETURNING *;

-- name: FindRefreshToken :one
SELECT * FROM refresh_tokens WHERE token_hash = $1;

-- name: InvalidRefreshToken :exec
UPDATE refresh_tokens 
SET expires_at = NOW() 
WHERE token_hash = $1;