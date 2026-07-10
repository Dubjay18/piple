-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  password,
  role
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, first_name, last_name, email, password, role, created_at;

-- name: GetUserByID :one
SELECT id, first_name, last_name, email, password, role, created_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password, role, created_at
FROM users
WHERE email = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2
WHERE id = $1
RETURNING id, first_name, last_name, email, password, role, created_at;