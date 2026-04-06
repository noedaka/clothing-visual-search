package minioclient

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/noedaka/clothing-visual-search/backend/internal/config"
	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type minioClient struct {
	client *minio.Client
	cfg    *config.Config
}

func NewMinIOClient(cfg *config.Config) (*minioClient, error) {
	client, err := minio.New(cfg.MinIOAddr, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	err = ensureMinIOBucket(client, cfg.MinIOBucket)
	if err != nil {
		return nil, err
	}

	log.Println("MinIO client initialized successfully")
	return &minioClient{client: client, cfg: cfg}, nil
}

func (c *minioClient) PutObject(ctx context.Context, objectKey string, imgData *model.ImageData) error {
	_, err := c.client.PutObject(ctx, c.cfg.MinIOBucket,
		objectKey, imgData.File, imgData.FileSize,
		minio.PutObjectOptions{ContentType: imgData.ContentType})

	return err
}

func (c *minioClient) RemoveObject(ctx context.Context, objectKey string) error {
	return c.client.RemoveObject(ctx, c.cfg.MinIOBucket, objectKey, minio.RemoveObjectOptions{})
}

func (c *minioClient) GenerateObjectKey(productID int64, filename string) string {
	name := filepath.Base(filename)
	uid := uuid.New().String()

	return fmt.Sprintf("products/%d/%s_%s", productID, uid, name)
}

func (c *minioClient) GetPublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", c.cfg.MinIOExternalBaseURL, c.cfg.MinIOBucket, objectKey)
}

func ensureMinIOBucket(client *minio.Client, minIOBucket string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, minIOBucket)
	if err != nil {
		return err
	}

	if !exists {
		if err := client.MakeBucket(ctx, minIOBucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		log.Printf("Bucket '%s' created successfully", minIOBucket)
	} else {
		log.Printf("Bucket '%s' already exists", minIOBucket)
	}

	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Principal": {"AWS": ["*"]},
            "Action": ["s3:GetObject"],
            "Resource": ["arn:aws:s3:::%s/*"]
        }]
    }`, minIOBucket)

	if err := client.SetBucketPolicy(ctx, minIOBucket, policy); err != nil {
		return fmt.Errorf("failed to set public policy: %w", err)
	}
	log.Printf("Public policy applied to bucket '%s'", minIOBucket)

	return nil
}
