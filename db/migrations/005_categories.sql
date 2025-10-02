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


-- +goose Down
DROP TABLE categories;