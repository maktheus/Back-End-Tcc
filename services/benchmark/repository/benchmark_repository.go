package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// BenchmarkRepository stores benchmark definitions.
type BenchmarkRepository struct {
	store *storage.MemoryRepository[models.Benchmark]
}

// New creates repo.
func New(store *storage.MemoryRepository[models.Benchmark]) *BenchmarkRepository {
	return &BenchmarkRepository{store: store}
}

// Save persists a benchmark.
func (r *BenchmarkRepository) Save(benchmark models.Benchmark) {
	r.store.Save(benchmark.ID, benchmark)
}

// List returns benchmarks.
func (r *BenchmarkRepository) List() []models.Benchmark {
	return r.store.List()
}
