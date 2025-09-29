package main

import (
	"go-twt/internal/store"
	"log"
	"net/http"
)

// Show login/signup page
func (app *application) authPageHandler(w http.ResponseWriter, r *http.Request) {
	data := store.PageData{
		CurrentUser: nil,
		Title:       "Welcome to go-twt",
	}

	app.render(w, "layout.html", data, "layout.html", "pages/auth.html")
}

// Handle login POST
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get user from database
	user, err := app.store.Users.GetByEmail(r.Context(), email)
	if err != nil {
		log.Printf("Login error: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// TODO: Verify password (we'll add bcrypt later)
	// For now, just check if password matches (INSECURE - fix later)
	if user.Password != password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// TODO: Create session (we'll add later)
	// For now, just redirect
	log.Printf("User logged in: %s", user.Email)

	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

// Handle signup POST
func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate
	if username == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// TODO: Hash password (we'll add bcrypt later)
	// For now, store plaintext (INSECURE - fix later)

	// Create user
	user := &store.User{
		Username: username,
		Email:    email,
		Password: password, // TODO: Hash this!
	}

	err := app.store.Users.Create(r.Context(), user)
	if err != nil {
		log.Printf("Signup error: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	log.Printf("User created: %s", user.Email)

	// TODO: Create session and redirect
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

// Handle logout
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Destroy session
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
