-- +goose Up
CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id),
  category TEXT NOT NULL CHECK (category IN ('lost_found','for_sale','event','general')),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  image_url TEXT[] DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ
);

-- +goose Down

Drop TABLE posts;