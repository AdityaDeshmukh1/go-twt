-- +goose Up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Optional indexes
CREATE INDEX idx_posts_user_id_created_at ON posts(user_id, created_at DESC);

-- +goose Down
DROP TABLE posts;
