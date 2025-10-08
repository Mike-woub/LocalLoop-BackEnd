-- name: LikePost :one
INSERT INTO likes (user_id, post_id) VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: UnlikePost :one
DELETE FROM likes WHERE user_id = $1 AND post_id = $2
RETURNING id;

-- name: GetLikeCount :one
SELECT COUNT(*) FROM likes WHERE post_id = $1;

-- name: CheckLiked :one
SELECT EXISTS (
  SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2
);