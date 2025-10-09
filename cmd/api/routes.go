package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Serve static assets (CSS/JS)
	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes (no auth required)
	r.Get("/", app.authPageHandler)
	r.Post("/v1/login", app.loginHandler)
	r.Post("/v1/signup", app.signupHandler)

	r.Group(func(r chi.Router) {
		r.Use(app.requireLogin)
		r.Get("/feed", app.feedHandler)
		r.Post("/tweet", app.createTweetHandler)
		r.Post("/logout", app.logoutHandler)
		r.Post("/like", app.toggleLikeHandler)

	})

	// API routes
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	log.Printf("Server started at %s", app.config.addr)
	return srv.ListenAndServe()
}
