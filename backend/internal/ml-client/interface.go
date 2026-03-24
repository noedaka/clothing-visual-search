package mlclient

import "context"

type MLClient interface {
	GetEmbedding(
		ctx context.Context,
		image []byte,
		imageFormat string,
	) ([]float32, error)
	Close() error
}
