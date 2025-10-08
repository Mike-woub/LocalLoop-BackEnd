-- +goose Up
CREATE TABLE likes (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  post_id INTEGER REFERENCES posts(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, post_id)
);

ALTER TABLE posts ADD COLUMN like_count INTEGER DEFAULT 0;


-- +goose Down

Drop TABLE likes;
