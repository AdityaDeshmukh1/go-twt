package store

import "time"

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-" db:"password_hash"`
	CreatedAt string `json:"created_at"`
}

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_id"`
	UpdatedAt string    `json:"updated_id"`

	// Joined fields (when fetching with author info)
	Author     *User `json:"author,omitempty"`
	LikeCount  int   `json:"like_count"`
	ReplyCount int   `json:"reply_count"`
}

// Template data structures
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
