package main

import (
	"go-twt/internal/store"
	"net/http"
)

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user := &store.User{
		Username: username,
		Email:    email,
		Password: password, // TODO: hash in production
	}

	err := app.store.Users.Create(r.Context(), user)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Signed up successfully!"))
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := app.store.Users.GetByEmail(r.Context(), email)
	if err != nil || user == nil {
		w.Write([]byte("Invalid email or password"))
		return
	}

	// TODO: compare hashed password in production
	if user.Password != password {
		w.Write([]byte("Invalid email or password"))
		return
	}

	w.Write([]byte("Welcome back, " + user.Username + "!"))
}
