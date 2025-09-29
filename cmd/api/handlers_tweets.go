package main

import (
	"go-twt/internal/store"
	"net/http"
)

// Show feed/timeline
func (app *application) feedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get current user from session
	// For now, use a mock user
	mockUser := &store.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Get posts from database
	posts, err := app.store.Posts.GetFeed(r.Context(), 20, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := store.FeedPageData{
		PageData: store.PageData{
			CurrentUser: mockUser,
			ActivePage:  "feed",
			Title:       "Home",
		},
		Posts:    posts,
		NextPage: 2,
	}

	app.render(w, "layout.html", data,
		"layout.html",
		"pages/feed.html",
		"components/composer.html",
		"components/tweet.html",
	)
}

// Create a new tweet
func (app *application) createTweetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get current user from session
	mockUser := &store.User{
		ID:       1,
		Username: "testuser",
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	post := &store.Post{
		UserID:  mockUser.ID,
		Content: content,
	}

	err := app.store.Posts.Create(r.Context(), post)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect back to feed
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}
