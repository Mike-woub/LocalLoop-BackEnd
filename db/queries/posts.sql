-- name: CreatePost :one
INSERT INTO posts (
  user_id, category_id, title, content, image_url, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPosts :many
SELECT 
  posts.id,
  posts.user_id,
  users.username,
  posts.title,
  posts.content,
  posts.image_url,
  posts.category_id,
  categories.name AS category_name
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
  categories.name AS category_name
FROM posts
JOIN categories ON posts.category_id = categories.id
JOIN users ON posts.user_id = users.id
WHERE posts.id=$1;
