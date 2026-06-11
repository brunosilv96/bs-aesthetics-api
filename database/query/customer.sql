
-- name: LoadCustomers :many
SELECT * FROM customers;

-- name: CreateCustomer :one
INSERT INTO customers (name, email, phone, password, birth_date, created_at)
VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING *;