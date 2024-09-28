-- name: CreateUser :execresult
INSERT INTO users(
    username, email, password
) VALUES (?, ?, ?);

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY id;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: UpdateUser :execresult
UPDATE users SET username = ?, email = ?, password = ? 
WHERE id=?;

-- name: DeleteUser :execresult
DELETE FROM users WHERE id = ?;
