-- name: CreatePost :exec
INSERT INTO posts(
    id, user_id, description, price
) VALUES (?, ?, ?, ?);

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY updated_at DESC, id ASC
LIMIT ? OFFSET ?;

-- name: GetPostById :one
SELECT * FROM posts
WHERE id = ? LIMIT 1;

-- name: UpdatePost :exec
UPDATE posts 
SET 
    description = COALESCE(sqlc.narg(description), description), 
    price = COALESCE(sqlc.narg(price), price), 
    sold = COALESCE(sqlc.narg(sold), sold)
WHERE id=?;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = ?;

