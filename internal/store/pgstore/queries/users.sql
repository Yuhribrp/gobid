-- name: CreateUser :one
INSERT INTO users (user_name, email, password_hash, bio)
VALUES ($1, $2, $3, $4)
RETURNING id, user_name, email, bio, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, user_name, email, bio, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, user_name, email, bio, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET user_name = COALESCE($2, user_name),
    email = COALESCE($3, email),
    bio = COALESCE($4, bio),
    updated_at = now()
WHERE id = $1
RETURNING id, user_name, email, bio, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, user_name, email, bio, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
