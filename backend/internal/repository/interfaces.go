package repository

import (
	"context"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type CategoryRepository interface {
	Add(ctx context.Context, category string) error
	List(ctx context.Context) ([]model.Category, error)
}

type ProductRepository interface {
	Add(ctx context.Context, product *model.Product) (int64, error)
	GetByIDs(ctx context.Context, IDs []int64) ([]model.Product, error)
}

type ImageRepository interface {
	Add(ctx context.Context, productID int64, imageData *model.ImageData) error
	GetByIDs(ctx context.Context, IDs []int64) ([]model.Image, error)
}
