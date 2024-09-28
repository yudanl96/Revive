-- name: CreatePost :execresult
INSERT INTO posts(
    user_id, description
) VALUES (?, ?);

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY updated_at;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = ? LIMIT 1;

-- name: UpdatePost :execresult
UPDATE posts SET description = ?
WHERE id=?;

-- name: DeletePost :execresult
DELETE FROM posts WHERE id = ?;

