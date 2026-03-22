package service

import (
	"context"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type CategoryService interface {
	Add(ctx context.Context, category string) error
	List(ctx context.Context) ([]model.Category, error)
}

type ProductService interface {
	Add(ctx context.Context, product *model.ProductWithImagesData) error
	GetByIDs(ctx context.Context, IDs []int64) ([]model.ProductWithImages, error)
}
