
-- name: LoadCustomers :many
SELECT * FROM customers;

-- name: CreateCustomer :one
INSERT INTO customers (id, name, email, phone, birth_date, created_at, updated_at, deleted_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;