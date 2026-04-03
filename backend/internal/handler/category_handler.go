package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

func (h *Handler) AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category *model.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.categoryService.Add(r.Context(), category.Name); err != nil {
		http.Error(w, "failed to add category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.List(r.Context())
	if err != nil {
		if errors.Is(err, model.ErrNoContent) {
			http.Error(w, "no categories", http.StatusNoContent)
			return
		}

		http.Error(w, "failed to list categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
