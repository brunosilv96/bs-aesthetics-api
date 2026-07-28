-- name: RegisterEnterprise :one
INSERT INTO enterprise (trade_name, cnpj, opening_time, closing_time, updated_at) 
VALUES ($1, $2, $3, $4, NULL) RETURNING *;