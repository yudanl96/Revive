-- name: CreateSavedPost :exec
INSERT INTO saved_posts(
    user_id, post_id
) VALUES (?, ?);

-- name: ListSavedPostsByUser :many
SELECT * FROM saved_posts
WHERE user_id = ?
ORDER BY saved_at DESC;

-- name: CountSavedPostsByPost :one
SELECT COUNT(*) FROM saved_posts
WHERE post_id = ?;

-- name: GetSavedPostByIds :one
SELECT * FROM saved_posts
WHERE user_id = ? AND post_id = ? LIMIT 1;


-- name: DeleteSavedPost :exec
DELETE FROM saved_posts 
WHERE user_id = ? AND post_id = ?;
