-- name: CreateComment :one
INSERT INTO comments (post_id, user_id, content, created_at ) VALUES(
    $1,$2,$3,$4
) RETURNING *;

-- name: GetComments :many
SELECT 
  c.content, 
  c.created_at, 
  u.username
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.post_id = $1;
