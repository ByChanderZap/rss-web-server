-- name: CreatePost :one
INSERT INTO posts (
  id,
  created_at,
  updated_at,
  title,
  description,
  published_at,
  url,
  feed_id
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.* 
FROM posts p
JOIN feeds_follows ff ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;