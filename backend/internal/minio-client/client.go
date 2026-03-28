package minioclient

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/noedaka/clothing-visual-search/backend/internal/config"
)

type minioClient struct {
	Client *minio.Client
	cfg *config.Config
}

func NewMinIOClient(cfg *config.Config) (*minioClient, error) {
	client, err := minio.New(cfg.MinIOAddr, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	log.Println("MinIO client initialized successfully")
	return &minioClient{Client: client, cfg: cfg}, nil
}

func (c *minioClient) EnsureMinIOBucket() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := c.Client.BucketExists(ctx, c.cfg.MinIOBucket)
	if err != nil {
		return err
	}

	if !exists {
		if err := c.Client.MakeBucket(ctx, c.cfg.MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		log.Printf("Bucket '%s' created successfully", c.cfg.MinIOBucket)
	} else {
		log.Printf("Bucket '%s' already exists", c.cfg.MinIOBucket)
	}

	return nil
}
