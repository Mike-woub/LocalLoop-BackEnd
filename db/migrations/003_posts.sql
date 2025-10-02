-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Seed default categories
INSERT INTO categories (name) VALUES
('General'),
('Lost & Found'),
('Events'),
('Jobs'),
('For Sale'),
('Services'),
('Recommendations'),
('Questions'),
('News & Alerts'),
('Rants & Raves'),
('Housing'),
('Transportation');

CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id),
  category_id INT NOT NULL REFERENCES categories(id),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  image_url TEXT[] DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE posts;
DROP TABLE categories;
