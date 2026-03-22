package milvusclient

import "context"

type MilvusInsertDelete interface {
	InsertEmbedding(ctx context.Context, imageID, productID int64, embedding []float32) error
	DeleteByImageID(ctx context.Context, imageID int64) error
}

type MilvusSearch interface {
	CreateCollection(ctx context.Context) error
	SearchSimilar(ctx context.Context, queryVector []float32, topK int) ([]int64, error)
}
