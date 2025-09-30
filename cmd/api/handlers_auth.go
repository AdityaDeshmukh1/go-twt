package main

import (
	"fmt"
	"go-twt/internal/store"
	"log"
	"net/http"
)

// Show login/signup page (uses auth-layout, NOT main layout)
func (app *application) authPageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Error string
	}{
		Error: r.URL.Query().Get("error"), // Get error from query params
	}

	// Use auth-layout.html instead of layout.html
	app.render(w, "auth-layout.html", data, "auth-layout.html", "pages/auth.html")
}

// Handle login POST
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate
	if email == "" || password == "" {
		http.Redirect(w, r, "/?error=All+fields+are+required", http.StatusSeeOther)
		return
	}

	// Get user from database
	user, err := app.store.Users.GetByEmail(r.Context(), email)
	if err != nil {
		log.Printf("Login error: %v", err)
		http.Redirect(w, r, "/?error=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	// TODO: Verify password with bcrypt (we'll add later)
	// For now, just check if password matches (INSECURE - fix later)
	if user.Password != password {
		http.Redirect(w, r, "/?error=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	session, _ := app.sessions.Get(r, "go-twt-session")
	session.Values["userID"] = user.ID
	session.Save(r, w)

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
		http.Redirect(w, r, "/?error=All+fields+are+required", http.StatusSeeOther)
		return
	}

	if len(password) < 8 {
		http.Redirect(w, r, "/?error=Password+must+be+at+least+8+characters", http.StatusSeeOther)
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
		http.Redirect(w, r, "/?error=Email+or+username+already+exists", http.StatusSeeOther)
		return
	}

	log.Printf("User created: %s", user.Email)

	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

// Handle logout
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessions.Get(r, "go-twt-session")
	session.Options.MaxAge = -1 // delete cookie
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) getCurrentUser(r *http.Request) (*store.User, error) {
	session, _ := app.sessions.Get(r, "go-twt-session")
	val, ok := session.Values["userID"].(int64)
	if !ok || val == 0 {
		return nil, fmt.Errorf("not logged in")
	}
	return app.store.Users.GetByID(r.Context(), val)
}

func (app *application) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.sessions.Get(r, "go-twt-session")
		userID, ok := session.Values["userID"].(int64)
		if !ok || userID == 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
