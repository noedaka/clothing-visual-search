package handler

import "github.com/noedaka/clothing-visual-search/backend/internal/service"

type Handler struct {
	productService   service.ProductService
	categoryService  service.CategoryService
	embeddingService service.EmbeddingService
	searchService    service.SearchService
}

func NewHandler(
	productService service.ProductService,
	categoryService service.CategoryService,
	embeddingService service.EmbeddingService,
	searchService service.SearchService,
) *Handler {
	return &Handler{
		productService:   productService,
		categoryService:  categoryService,
		embeddingService: embeddingService,
		searchService:    searchService,
	}
}
