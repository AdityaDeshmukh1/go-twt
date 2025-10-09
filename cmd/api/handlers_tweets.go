package main

import (
	"go-twt/internal/store"
	"net/http"
	"strconv"
)

// Show feed/timeline
func (app *application) feedHandler(w http.ResponseWriter, r *http.Request) {

	user, err := app.getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get posts from database
	posts, err := app.store.Posts.GetFeed(r.Context(), 20, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := store.FeedPageData{
		PageData: store.PageData{
			CurrentUser: user,
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

func (app *application) toggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
	}

	postIDStr := r.FormValue("post_id")
	if postIDStr == "" {
		http.Error(w, "Missing post_id", http.StatusBadRequest)
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	like := &store.Like{
		UserID: user.ID,
		PostID: postID,
	}

	exists, err := app.store.Likes.Exists(r.Context(), like)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if exists {
		err = app.store.Likes.Delete(r.Context(), like)
	} else {
		err = app.store.Likes.Create(r.Context(), like)
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect back to the feed (or post page)
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}
