-- name: CreateSavedPost :execresult
INSERT INTO saved_posts(
    user_id, post_id
) VALUES (?, ?);

-- name: ListSavedPostsByUser :many
SELECT * FROM saved_posts
WHERE user_id = ?
ORDER BY saved_at DESC;

-- name: ListSavedPostsByPost :many
SELECT * FROM saved_posts
WHERE post_id = ?
ORDER BY saved_at DESC;

-- name: GetSavedPostByUser :many
SELECT * FROM saved_posts
WHERE user_id = ? LIMIT ? OFFSET ?;

-- name: GetSavedPostByPost :many
SELECT * FROM saved_posts
WHERE post_id = ? LIMIT ? OFFSET ?;


-- name: DeleteSavedPost :execresult
DELETE FROM saved_posts 
WHERE user_id = ? AND post_id = ?;
