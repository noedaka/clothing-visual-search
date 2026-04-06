package handler

import (
	"github.com/noedaka/clothing-visual-search/backend/internal/config"
	"github.com/noedaka/clothing-visual-search/backend/internal/service"
)

type Handler struct {
	productService   service.ProductService
	categoryService  service.CategoryService
	embeddingService service.EmbeddingService
	searchService    service.SearchService
	cfg              config.Config
}

func NewHandler(
	productService service.ProductService,
	categoryService service.CategoryService,
	embeddingService service.EmbeddingService,
	searchService service.SearchService,
	cfg config.Config,
) *Handler {
	return &Handler{
		productService:   productService,
		categoryService:  categoryService,
		embeddingService: embeddingService,
		searchService:    searchService,
		cfg:              cfg,
	}
}
