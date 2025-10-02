package store

import (
	"context"
	"database/sql"
)

type RetweetStore struct {
	db *sql.DB
}

// Create inserts a new retweet (or quote-tweet)
func (s *RetweetStore) Create(ctx context.Context, retweet *Retweet) error {
	query := `
	INSERT INTO retweets (user_id, post_id, comment)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, post_id) DO NOTHING
	RETURNING id, created_at
	`

	return s.db.QueryRowContext(ctx, query, retweet.UserID, retweet.PostID, retweet.Comment).
		Scan(&retweet.ID, &retweet.CreatedAt)
}

// Delete removes a retweet by user and post
func (s *RetweetStore) Delete(ctx context.Context, userID, postID int64) error {
	query := `
	DELETE FROM retweets
	WHERE user_id = $1 AND post_id = $2
	`

	_, err := s.db.ExecContext(ctx, query, userID, postID)
	return err
}

// Exists checks if the user has already retweeted the post
func (s *RetweetStore) Exists(ctx context.Context, userID, postID int64) (bool, error) {
	query := `
	SELECT 1
	FROM retweets
	WHERE user_id = $1 AND post_id = $2
	LIMIT 1
	`

	var tmp int
	err := s.db.QueryRowContext(ctx, query, userID, postID).Scan(&tmp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// CountByPost returns total number of retweets for a post
func (s *RetweetStore) CountByPost(ctx context.Context, postID int64) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM retweets
	WHERE post_id = $1
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
