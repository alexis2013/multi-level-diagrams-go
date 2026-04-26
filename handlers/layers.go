package handlers

import (
	"html/template"
	"layerdraw/models"
	"layerdraw/storage"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

type LayerHandler struct {
	Store storage.LayerStore
}

func (h *LayerHandler) List(w http.ResponseWriter, r *http.Request) {
	layers := h.Store.List()
	renderLayers(w, layers)
}

func (h *LayerHandler) Add(w http.ResponseWriter, r *http.Request) {
	h.Store.Add("New Layer")
	layers := h.Store.List()
	renderLayers(w, layers)
}

func (h *LayerHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := h.Store.Toggle(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	layers := h.Store.List()
	renderLayers(w, layers)
}

func (h *LayerHandler) Rename(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	if _, err := h.Store.Rename(id, name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	layers := h.Store.List()
	renderLayers(w, layers)
}

func renderLayers(w http.ResponseWriter, layers []models.Layer) {
	tmpl := template.Must(template.ParseFiles("templates/layer_panel.html", "templates/layer_item.html"))
	data := struct {
		Layers []models.Layer
	}{
		Layers: layers,
	}
	if err := tmpl.ExecuteTemplate(w, "layer_panel.html", data); err != nil {
		log.Printf("renderLayers: failed to execute template: %v", err)
	}
}
