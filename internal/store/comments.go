package store

import (
	"context"
	"database/sql"
)

type CommentStore struct {
	db *sql.DB
}

// Create inserts a new comment
func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
	INSERT INTO comments (user_id, post_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err := s.db.QueryRowContext(ctx, query, comment.UserID, comment.PostID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)

	return err
}

// Delete removes a comment by commentID and userID (only author can delete)
func (s *CommentStore) Delete(ctx context.Context, commentID, userID int64) error {
	query := `
	DELETE FROM comments
	WHERE id = $1 AND user_id = $2
	`

	_, err := s.db.ExecContext(ctx, query, commentID, userID)
	return err
}

// GetByPost fetches comments for a post with pagination
func (s *CommentStore) GetByPost(ctx context.Context, postID int64, limit, offset int) ([]Comment, error) {
	query := `
	SELECT id, user_id, post_id, content, created_at
	FROM comments
	WHERE post_id = $1
	ORDER BY created_at ASC
	LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.PostID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// CountByPost returns total number of comments for a post
func (s *CommentStore) CountByPost(ctx context.Context, postID int64) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM comments
	WHERE post_id = $1
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
