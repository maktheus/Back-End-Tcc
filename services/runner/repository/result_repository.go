package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// ResultRepository stores submission results.
type ResultRepository struct {
	store *storage.MemoryRepository[models.Submission]
}

// New creates repository.
func New(store *storage.MemoryRepository[models.Submission]) *ResultRepository {
	return &ResultRepository{store: store}
}

// Save stores submission.
func (r *ResultRepository) Save(sub models.Submission) {
	r.store.Save(sub.ID, sub)
}

// List returns submissions.
func (r *ResultRepository) List() []models.Submission {
	return r.store.List()
}
