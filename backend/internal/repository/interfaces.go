package repository

import (
	"context"
	"io"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type CategoryRepository interface {
	Add(ctx context.Context, category string) error
	List(ctx context.Context) ([]model.Category, error)
}

type ProductRepository interface {
	Add(ctx context.Context, product *model.Product) (int, error)
	GetByIDs(ctx context.Context, ids []int) ([]model.Product, error)
}

type ImageRepository interface {
	Add(ctx context.Context, productID int, file io.Reader,
		fileSize int64, filename, contentType string, isPrimary bool) error
	GetByIDs(ctx context.Context, productIDs []int) ([]model.Image, error)
}
