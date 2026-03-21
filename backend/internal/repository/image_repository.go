package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/noedaka/clothing-visual-search/backend/internal/config"
	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type ImageRepo struct {
	db          *sql.DB
	minioClient *minio.Client
	cfg         *config.Config
}

func NewImageRepo(db *sql.DB, minioClient *minio.Client, cfg *config.Config) *ImageRepo {
	return &ImageRepo{db: db, minioClient: minioClient, cfg: cfg}
}

func (r ImageRepo) generateObjectKey(productID int, filename string) string {
	name := filepath.Base(filename)
	uid := uuid.New().String()

	return fmt.Sprintf("products/%d/%s_%s", productID, uid, name)
}

func (r *ImageRepo) getPublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", r.cfg.MinIOExternalBaseURL, r.cfg.MinIOBucket, objectKey)
}

func (r *ImageRepo) Add(
	ctx context.Context,
	productID int,
	imageData *model.ImageData,
) error {
	objectKey := r.generateObjectKey(productID, imageData.Filename)

	_, err := r.minioClient.PutObject(ctx, r.cfg.MinIOBucket, objectKey, imageData.File, imageData.FileSize,
		minio.PutObjectOptions{
			ContentType: imageData.ContentType,
		})

	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			err = r.minioClient.RemoveObject(
				ctx, r.cfg.MinIOBucket, objectKey, minio.RemoveObjectOptions{})
			if err != nil {
				log.Printf("failed to remove object after error: %v", err)
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO product_images (product_id, object_key, is_primary)
		VALUES ($1, $2, $3)`,
	)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ImageRepo) GetByIDs(ctx context.Context, IDs []int) ([]model.Image, error) {
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
