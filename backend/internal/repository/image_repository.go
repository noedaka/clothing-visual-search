package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/noedaka/clothing-visual-search/backend/internal/config"
	milvusclient "github.com/noedaka/clothing-visual-search/backend/internal/milvus-client"
	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type ImageRepo struct {
	db           *sql.DB
	minioClient  *minio.Client
	milvusClient milvusclient.MilvusInsertDelete
	cfg          *config.Config
}

func NewImageRepo(
	db *sql.DB,
	minioClient *minio.Client,
	milvusClient milvusclient.MilvusInsertDelete,
	cfg *config.Config,
) *ImageRepo {
	return &ImageRepo{db: db, minioClient: minioClient, milvusClient: milvusClient, cfg: cfg}
}

func (r *ImageRepo) Add(
	ctx context.Context,
	tx *sql.Tx,
	productID int64,
	imageData *model.ImageData,
) (int64, string, error) {
	objectKey := r.generateObjectKey(productID, imageData.Filename)

	_, err := r.minioClient.PutObject(ctx, r.cfg.MinIOBucket, objectKey,
		imageData.File, imageData.FileSize,
		minio.PutObjectOptions{ContentType: imageData.ContentType})
	if err != nil {
		return 0, "", err
	}

	var imageID int64
	milvusInserted := false
	defer func() {
		if err != nil {
			_ = r.minioClient.RemoveObject(ctx, r.cfg.MinIOBucket, objectKey, minio.RemoveObjectOptions{})
			if milvusInserted {
				_ = r.milvusClient.DeleteByImageID(ctx, imageID)
			}
		}
	}()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO product_images (product_id, object_key, is_primary)
         VALUES ($1, $2, $3) RETURNING id`,
		productID, objectKey, imageData.IsPrimary).Scan(&imageID)
	if err != nil {
		return 0, "", err
	}

	err = r.milvusClient.InsertEmbedding(ctx, imageID, productID, imageData.Embedding)
	if err != nil {
		return 0, "", err
	}
	milvusInserted = true

	return imageID, objectKey, nil
}

func (r *ImageRepo) GetByIDs(ctx context.Context, IDs []int64) ([]model.Image, error) {
	if len(IDs) == 0 {
		return []model.Image{}, nil
	}

	query := `
		SELECT id, product_id, object_key, is_primary
		FROM product_images
		WHERE product_id = ANY($1) ORDER BY product_id, is_primary DESC
	`
	args := []interface{}{IDs}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoContent
		}
		return nil, err
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		if err := rows.Scan(&img.ID, &img.ProductID, &img.ObjectKey, &img.IsPrimary); err != nil {
			return nil, err
		}

		img.URL = r.getPublicURL(img.ObjectKey)
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepo) DeleteByID(ctx context.Context, ID int64, objectKey string) error {
	_ = r.minioClient.RemoveObject(ctx, r.cfg.MinIOBucket, objectKey, minio.RemoveObjectOptions{})
	_ = r.milvusClient.DeleteByImageID(ctx, ID)

	return nil
}

func (r ImageRepo) generateObjectKey(productID int64, filename string) string {
	name := filepath.Base(filename)
	uid := uuid.New().String()

	return fmt.Sprintf("products/%d/%s_%s", productID, uid, name)
}

func (r *ImageRepo) getPublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", r.cfg.MinIOExternalBaseURL, r.cfg.MinIOBucket, objectKey)
}
