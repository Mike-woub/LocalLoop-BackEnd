-- name: CreatePost :one
INSERT INTO posts (
  user_id, category, title, content, image_url, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;
