package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

func (h *Handler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "file too large or invalid format", http.StatusBadRequest)
		return
	}

	productJSON := r.FormValue("product")
	if productJSON == "" {
		http.Error(w, "missing prod data", http.StatusBadRequest)
		return
	}

	var product model.Product
	if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
		http.Error(w, "invalid JSON"+err.Error(), http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if product.Price <= 0 {
		http.Error(w, "price must be positive", http.StatusBadRequest)
		return
	}
	if product.CategoryID <= 0 {
		http.Error(w, "category_id is required", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "no images", http.StatusBadRequest)
		return
	}

	var imagesData []model.ImageData
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "failed to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		imagesByte, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(fileHeader.Filename)
		format := strings.TrimPrefix(ext, ".")
		if format == "" {
			format = "jpg"
		}

		embedding, err := h.embeddingService.GetEmbedding(r.Context(), imagesByte, format)
		if err != nil {
			http.Error(w, "error getting embedding", http.StatusInternalServerError)
			return
		}

		imgData := model.ImageData{
			File:        bytes.NewReader(imagesByte),
			FileSize:    fileHeader.Size,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Contetn-Type"),
			IsPrimary:   i == 0,
			Embedding:   embedding,
		}

		imagesData = append(imagesData, imgData)
	}

	req := model.ProductWithImagesData{
		Product:           product,
		ProductImagesData: imagesData,
	}

	err := h.productService.Add(r.Context(), &req)
	if err != nil {
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SearchByImageHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read image", http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	format := strings.TrimPrefix(ext, ".")
	if format == "" {
		format = "jpg"
	}

	embedding, err := h.embeddingService.GetEmbedding(r.Context(), imageData, format)
	if err != nil {
		http.Error(w, "ml service error", http.StatusInternalServerError)
		return
	}

	topK := 10
	productIDs, err := h.searchService.SearchSimilar(r.Context(), embedding, topK)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	if len(productIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	productsWithImages, err := h.productService.GetByIDs(r.Context(), productIDs)
	if err != nil {
		if errors.Is(err, model.ErrNoContent) {
			http.Error(w, "no prodcuts", http.StatusNoContent)
			return
		}

		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(productsWithImages); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
