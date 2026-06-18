
-- name: LoadCustomers :many
SELECT * FROM customers;

-- name: CreateCustomer :one
INSERT INTO customers (name, email, phone, password, birthdate, created_at)
VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING *;

-- name: FindCustomerByID :one
SELECT * FROM customers WHERE id = $1;

-- name: SoftDeleteCustomer :exec
UPDATE customers SET updated_at = NOW(), deleted_at = NOW() WHERE id = $1;

-- name: UpdateCustomer :exec
UPDATE customers 
SET name = $2, email = $3, phone = $4, birthdate = $5, updated_at = NOW() 
WHERE id = $1;