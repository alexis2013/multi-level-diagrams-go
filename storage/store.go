package storage

import (
	"fmt"
	"layerdraw/models"
	"sync"
	"time"
)

type LayerStore interface {
	List() []models.Layer
	Add(name string) models.Layer
	Toggle(id string) (models.Layer, error)
	Rename(id string, name string) (models.Layer, error)
}

type MemoryStore struct {
	mu     sync.Mutex
	layers []models.Layer
}

func NewMemoryStore() *MemoryStore {
	ms := &MemoryStore{
		layers: []models.Layer{},
	}
	ms.Add("Layer 1")
	return ms
}

func (s *MemoryStore) List() []models.Layer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]models.Layer(nil), s.layers...)
}

func (s *MemoryStore) Add(name string) models.Layer {
	s.mu.Lock()
	defer s.mu.Unlock()
	layer := models.Layer{
		ID:      fmt.Sprintf("layer-%d", time.Now().UnixNano()),
		Name:    name,
		Visible: true,
	}
	s.layers = append(s.layers, layer)
	return layer
}

func (s *MemoryStore) Toggle(id string) (models.Layer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, l := range s.layers {
		if l.ID == id {
			s.layers[i].Visible = !s.layers[i].Visible
			return s.layers[i], nil
		}
	}
	return models.Layer{}, fmt.Errorf("layer %s not found", id)
}

func (s *MemoryStore) Rename(id string, name string) (models.Layer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, l := range s.layers {
		if l.ID == id {
			s.layers[i].Name = name
			return s.layers[i], nil
		}
	}
	return models.Layer{}, fmt.Errorf("layer %s not found", id)
}
