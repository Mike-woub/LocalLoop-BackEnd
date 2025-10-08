-- +goose Up
ALTER TABLE posts
ADD COLUMN latitude DOUBLE PRECISION DEFAULT 9.03,
ADD COLUMN longitude DOUBLE PRECISION DEFAULT 38.74,
ADD COLUMN location_name VARCHAR(255) DEFAULT 'Addis Ababa, Ethiopia';

-- +goose Down
ALTER TABLE posts
DROP COLUMN latitude,
DROP COLUMN longitude,
DROP COLUMN location_name;
