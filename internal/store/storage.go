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
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}

}
