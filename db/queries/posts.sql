-- name: CreatePost :one
INSERT INTO posts (
  user_id,
  category_id,
  title,
  content,
  image_url,
  expires_at,
  latitude,
  longitude,
  location_name
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetPosts :many
SELECT 
  posts.id,
  posts.user_id,
  users.username,
  posts.title,
  posts.content,
  posts.image_url,
  posts.category_id,
  categories.name AS category_name,
  posts.latitude,
  posts.longitude,
  posts.location_name
FROM posts
JOIN categories ON posts.category_id = categories.id
JOIN users ON posts.user_id = users.id;

-- name: GetCertainPost :one
SELECT 
  posts.id,
  posts.user_id,
  users.username,
  posts.title,
  posts.content,
  posts.image_url,
  posts.category_id,
  categories.name AS category_name,
  posts.latitude,
  posts.longitude,
  posts.location_name
FROM posts
JOIN categories ON posts.category_id = categories.id
JOIN users ON posts.user_id = users.id
WHERE posts.id = $1;

-- name: DeletePost :one
DELETE FROM posts 
WHERE id = $1 AND user_id = $2
RETURNING id;
