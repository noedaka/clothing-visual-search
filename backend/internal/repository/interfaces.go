package repository

import (
	"context"
	"database/sql"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type CategoryRepository interface {
	Add(ctx context.Context, category string) error
	List(ctx context.Context) ([]model.Category, error)
}

type ProductRepository interface {
	Add(ctx context.Context, tx *sql.Tx, product *model.Product) (int64, error)
	GetByIDs(ctx context.Context, IDs []int64) ([]model.Product, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type ImageRepository interface {
	Add(
		ctx context.Context,
		tx *sql.Tx,
		productID int64,
		imageData *model.ImageData,
	) (int64, string, error)
	GetByIDs(ctx context.Context, IDs []int64) ([]model.Image, error)
	DeleteByID(ctx context.Context, ID int64, objectKey string) error
}
