package milvusclient

import "context"

type MilvusInsertDelete interface {
	InsertEmbedding(ctx context.Context, imageID, productID int64, embedding []float32) error
	DeleteByImageID(ctx context.Context, imageID int64) error
}

type MilvusSearcher interface {
	SearchSimilar(ctx context.Context, queryVector []float32, topK int, threshold float64) ([]int64, error)
}
