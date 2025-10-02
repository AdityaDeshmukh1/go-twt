package store

import (
	"context"
	"database/sql"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Content,
	).Scan(&post.ID, &post.CreatedAt)

	return err
}

func (s *PostStore) GetFeed(ctx context.Context, limit int, offset int) ([]Post, error) {
	query := `
	SELECT 
		p.id, 
		p.user_id, 
		p.content, 
		p.created_at,
		u.id,
		u.username,
		u.email,
		u.created_at,
		(SELECT COUNT(*) FROM likes WHERE post_id = p.id) AS like_count,
		(SELECT COUNT(*) FROM retweets WHERE post_id = p.id) AS retweet_count,
		(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS reply_count
	FROM posts p
	JOIN users u ON p.user_id = u.id
	ORDER BY p.created_at DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var user User

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&post.LikeCount,
			&post.RetweetCount,
			&post.ReplyCount,
		)
		if err != nil {
			return nil, err
		}

		post.Author = &user
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (s *PostStore) GetByUserID(ctx context.Context, userID int64, limit int) ([]Post, error) {
	query := `
		SELECT 
			p.id, 
			p.user_id, 
			p.content, 
			p.created_at,
			u.id,
			u.username,
			u.email,
			u.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var user User

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		post.Author = &user
		posts = append(posts, post)
	}

	return posts, rows.Err()
}
