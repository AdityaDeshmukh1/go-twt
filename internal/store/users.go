package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-" db:"password_hash"`
	CreatedAt string `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx,
		query,
		user.Username,
		user.Email,
		user.Password).
		Scan(
			&user.ID,
			&user.CreatedAt,
		)

	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	query := `
	SELECT id, username, email, password_hash, created_at
	FROM users
	WHERE email = $1`

	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}
