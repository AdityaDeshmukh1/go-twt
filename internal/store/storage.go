package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetFeed(ctx context.Context, limit int, offset int) ([]Post, error)
		GetByUserID(ctx context.Context, userID int64, limit int) ([]Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
		GetByUsername(ctx context.Context, username string) (*User, error)
		GetByID(ctx context.Context, id int64) (*User, error)
	}
	Likes interface {
		Create(ctx context.Context, like *Like) error
		Delete(ctx context.Context, like *Like) error
		Exists(ctx context.Context, like *Like) (bool, error)
		CountByPost(ctx context.Context, postID int64) (int, error)
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		Delete(ctx context.Context, commentID, userID int64) error
		GetByPost(ctx context.Context, postID int64, limit, offset int) ([]Comment, error)
		CountByPost(ctx context.Context, postID int64) (int, error)
	}
	Retweets interface {
		Create(ctx context.Context, retweet *Retweet) error
		Delete(ctx context.Context, userID, postID int64) error
		Exists(ctx context.Context, userID, postID int64) (bool, error)
		CountByPost(ctx context.Context, postID int64) (int, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Likes:    &LikeStore{db},
		Comments: &CommentStore{db},
		Retweets: &RetweetStore{db},
	}

}
