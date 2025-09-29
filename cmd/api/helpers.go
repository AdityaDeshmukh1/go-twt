package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// parseTemplate parses template files from web/templates
func (app *application) parseTemplate(files ...string) (*template.Template, error) {
	// Add template functions
	funcMap := template.FuncMap{
		"timeAgo": timeAgo,
	}

	// Prepend the full path to each file
	fullPaths := make([]string, len(files))
	for i, file := range files {
		fullPaths[i] = filepath.Join("web", "templates", file)
	}

	return template.New("").Funcs(funcMap).ParseFiles(fullPaths...)
}

// render executes a template with data
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

// serverError logs error and sends 500
func (app *application) serverError(w http.ResponseWriter, err error) {
	log.Printf("Server error: %v", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// clientError sends a specific status code
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound sends 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// timeAgo returns a human-readable time difference
func timeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}
	return t.Format("Jan 2, 2006")
}
