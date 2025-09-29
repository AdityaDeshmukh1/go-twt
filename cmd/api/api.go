package main

import (
	"go-twt/internal/store"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// application holds config, db store
type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// Helper function to parse templates for a specific page
func (app *application) parseTemplate(files ...string) (*template.Template, error) {
	// Prepend the full path to each file
	fullPaths := make([]string, len(files))
	for i, file := range files {
		fullPaths[i] = filepath.Join("web", "templates", file)
	}
	return template.ParseFiles(fullPaths...)
}

// Helper function to render a template with error handling
func (app *application) render(w http.ResponseWriter, templateName string, data interface{}, files ...string) {
	tmpl, err := app.parseTemplate(files...)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

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

	// Frontend routes
	r.Get("/", app.indexHandler)
	r.Get("/feed", app.feedHandler)

	// API routes
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Post("/signup", app.signupHandler)
		r.Post("/login", app.loginHandler)
	})

	return r
}

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse layout.html and index.html together
	app.render(w, "layout.html", nil, "layout.html", "index.html")
}

func (app *application) feedHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Posts []string
	}{
		Posts: []string{"First post!", "Hello world!", "Go templates are working!"},
	}

	// Parse layout.html and feed.html together
	app.render(w, "layout.html", data, "layout.html", "feed.html")
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
