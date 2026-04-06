package minioclient

import (
	"context"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type MinIOClient interface {
	PutObject(ctx context.Context, objectKey string, imgData *model.ImageData) error
	RemoveObject(ctx context.Context, objectKey string) error
	GenerateObjectKey(productID int64, filename string) string
	GetPublicURL(objectKey string) string
}