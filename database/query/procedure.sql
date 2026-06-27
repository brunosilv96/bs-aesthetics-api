-- name: CreateProcedure :one
INSERT INTO procedures (registred_by, name, description, banner_url, price, duration_minutes, available, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: FindProcedureByID :one
SELECT * FROM procedures WHERE id = $1;