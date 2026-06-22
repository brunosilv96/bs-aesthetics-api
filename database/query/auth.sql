-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (customer_id, token_hash, expires_at, created_at)
VALUES ($1, $2, $3, NOW()) RETURNING *;