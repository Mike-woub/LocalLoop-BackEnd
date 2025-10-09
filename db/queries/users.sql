-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES($1,$2,$3)
RETURNING id, username, email, created_at;

-- name: GetUserByEmail :one
SELECT id, username, email, password, created_at, avatar_url
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, username, email, password, avatar_url
FROM users
WHERE id = $1;


-- name: UpdateUsername :exec
UPDATE users
SET username = $1
WHERE id = $2;

-- name: UpdateEmail :exec
UPDATE users
SET email = $1
WHERE id = $2;

-- name: UpdatePassword :exec
UPDATE users
SET password = $1
WHERE id = $2;

-- name: UpdateAvatar :exec
UPDATE users
SET avatar_url = $1
WHERE id = $2;

-- name: CheckUsernameExists :one
SELECT COUNT(*) FROM users WHERE username = $1 AND id != $2;

-- name: CheckEmailExists :one
SELECT COUNT(*) FROM users WHERE email = $1 AND id != $2;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE email = $1;
