-- name: CreateUser :exec
INSERT INTO users(
    id, username, email, password
) VALUES (UUID(),?, ?, ?);

-- name: RetrieveIdByEmail :one
SELECT id FROM users
WHERE email = ? LIMIT 1;

-- name: RetrieveIdByUsername :one
SELECT id FROM users
WHERE username = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY id 
LIMIT ? OFFSET ?;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users SET username = ?, email = ?, password = ?
WHERE id=?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
