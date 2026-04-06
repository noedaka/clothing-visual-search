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

type EmbeddingService interface {
	GetEmbedding(
		ctx context.Context,
		image []byte,
		imageFormat string,
	) ([]float32, error)
}

type SearchService interface {
	SearchSimilar(
		ctx context.Context,
		queryVector []float32,
		topK int,
		threshold float64,
	) ([]int64, error)
}
