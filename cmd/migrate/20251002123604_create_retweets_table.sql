-- +goose Up
CREATE TABLE retweets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    comment TEXT, -- optional quote tweet
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_post UNIQUE(user_id, post_id)
);

-- Index for counting retweets by post
CREATE INDEX idx_retweets_post_id ON retweets(post_id);

-- +goose Down
DROP TABLE IF EXISTS retweets;
