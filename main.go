package main

import (
	"fmt"
	"html/template"
	"layerdraw/handlers"
	"layerdraw/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var templates *template.Template

func main() {
	// Parse templates at startup
	var err error
	templates, err = template.ParseFiles("templates/base.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// Initialize store
	store := storage.NewMemoryStore()

	// Initialize handlers
	layerHandler := &handlers.LayerHandler{Store: store}

	r := chi.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Routes
	r.Get("/", indexHandler)
	r.Get("/layers", layerHandler.List)
	r.Post("/layers", layerHandler.Add)
	r.Post("/layers/{id}/toggle", layerHandler.Toggle)
	r.Post("/layers/{id}/rename", layerHandler.Rename)

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(fmt.Errorf("server failed: %w", err))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := templates.ExecuteTemplate(w, "base.html", nil); err != nil {
		log.Printf("failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
