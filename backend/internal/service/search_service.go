package service

import (
	"context"

	milvusclient "github.com/noedaka/clothing-visual-search/backend/internal/milvus-client"
)

type SearchServ struct {
	client milvusclient.MilvusSearcher
}

func NewSearchServ(client milvusclient.MilvusSearcher) *SearchServ {
	return &SearchServ{client: client}
}

func (s *SearchServ) SearchSimilar(ctx context.Context, queryVector []float32, topK int) ([]int64, error) {
	return s.client.SearchSimilar(ctx, queryVector, topK)
}
