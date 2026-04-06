package service

import (
	"context"

	mlclient "github.com/noedaka/clothing-visual-search/backend/internal/ml-client"
)

type EmbeddingServ struct {
	client mlclient.MLClient
}

func NewEmbeddingServ(client mlclient.MLClient) *EmbeddingServ {
	return &EmbeddingServ{client: client}
}

func (s *EmbeddingServ) GetEmbedding(
	ctx context.Context,
	image []byte,
	imageFormat string,
) ([]float32, error) {
	return s.client.GetEmbedding(ctx, image, imageFormat)
}

