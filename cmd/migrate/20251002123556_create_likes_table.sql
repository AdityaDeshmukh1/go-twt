-- +goose Up
CREATE TABLE likes (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

-- Index on post_id for fast counting
CREATE INDEX idx_likes_post_id ON likes(post_id);

-- +goose Down
DROP TABLE IF EXISTS likes;
