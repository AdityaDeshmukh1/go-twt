package store

import (
	"context"
	"database/sql"
)

type LikeStore struct {
	db *sql.DB
}

func (s *LikeStore) Create(ctx context.Context, like *Like) error {
	query := `
	INSERT INTO likes (user_id, post_id)
	VALUES ($1, $2)
    ON CONFLICT (user_id, post_id) DO NOTHING
	RETURNING created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		like.UserID,
		like.PostID,
	).Scan(&like.CreatedAt)

	return err
}

func (s *LikeStore) Delete(ctx context.Context, like *Like) error {
	query := `
	DELETE from likes 
	WHERE user_id = $1 AND post_id = $2 
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		like.UserID,
		like.PostID)

	return err
}

func (s *LikeStore) Exists(ctx context.Context, like *Like) (bool, error) {
	query := `
		SELECT 1
		FROM likes
		WHERE user_id = $1 AND post_id = $2
		LIMIT 1;
	`

	var tmp int
	err := s.db.QueryRowContext(ctx, query, like.UserID, like.PostID).Scan(&tmp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *LikeStore) CountByPost(ctx context.Context, postID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM likes
		WHERE post_id = $1;
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
