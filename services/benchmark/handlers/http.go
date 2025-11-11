package handlers

import (
	"encoding/json"
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/services/benchmark/service"
)

// HTTP handles benchmark endpoints.
type HTTP struct {
	service *service.Service
}

// New creates handlers.
func New(service *service.Service) *HTTP {
	return &HTTP{service: service}
}

// Create handles creation requests.
func (h *HTTP) Create(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkghttp.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	benchmark, err := h.service.Create(payload.ID, payload.Name, payload.Description)
	if err != nil {
		pkghttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	pkghttp.JSON(w, http.StatusCreated, benchmark)
}

// List returns benchmarks.
func (h *HTTP) List(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.List())
}
