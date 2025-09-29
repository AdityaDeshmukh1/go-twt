-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username     TEXT NOT NULL UNIQUE,
    email        TEXT NOT NULL UNIQUE,
    password     TEXT NOT NULL, -- store bcrypt hash
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;

