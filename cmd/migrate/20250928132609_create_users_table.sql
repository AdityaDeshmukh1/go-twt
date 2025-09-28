-- +goose Up
-- Enable citext extension if not already enabled
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email CITEXT NOT NULL UNIQUE, -- case-insensitive email
    password_hash VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;

