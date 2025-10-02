-- +goose Up
CREATE TABLE neighborhoods (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    lat NUMERIC NOT NULL,
    lng NUMERIC NOT NULL,
    radius_km NUMERIC DEFAULT 5
);

-- +goose Down
DROP TABLE neighborhoods;