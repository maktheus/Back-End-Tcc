package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// Repository stores trace events.
type Repository struct {
	store *storage.MemoryRepository[models.TraceEvent]
}

// New creates repository.
func New(store *storage.MemoryRepository[models.TraceEvent]) *Repository {
	return &Repository{store: store}
}

// Save persists event.
func (r *Repository) Save(event models.TraceEvent) {
	r.store.Save(event.ID, event)
}

// List returns events.
func (r *Repository) List() []models.TraceEvent {
	return r.store.List()
}
