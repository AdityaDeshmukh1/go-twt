package store

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Post struct {
	ID        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	Author       *User `json:"author,omitempty"`
	LikeCount    int   `json:"like_count"`
	RetweetCount int   `json:"retweet_count`
	ReplyCount   int   `json:"reply_count"`
}

type Like struct {
	UserID    int64     `json:"user_id" db:"user_id"`
	PostID    int64     `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	User *User `json:"user,omitempty"`
	Post *Post `json:"post,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	PostID    int64     `json:"post_id" db:"post_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	Author *User `json:"author,omitempty"`
	Post   *Post `json:"post,omitempty"`
}

type Retweet struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	PostID    int64     `json:"post_id" db:"post_id"`
	Comment   *string   `json:"comment,omitempty" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	User *User `json:"user,omitempty"`
	Post *Post `json:"post,omitempty"`
}

type PageData struct {
	CurrentUser *User
	ActivePage  string
	Title       string
	Error       string
}

type FeedPageData struct {
	PageData
	Posts    []Post
	NextPage int
	HasMore  bool
}

type ProfilePageData struct {
	PageData
	User           *User
	Posts          []Post
	IsOwnProfile   bool
	IsFollowing    bool
	FollowerCount  int
	FollowingCount int
}
