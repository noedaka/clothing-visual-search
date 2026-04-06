package service

import (
	"context"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
	"github.com/noedaka/clothing-visual-search/backend/internal/repository"
)

type ProductServ struct {
	productRepo repository.ProductRepository
	imageRepo   repository.ImageRepository
}

func NewProductServ(
	productRepo repository.ProductRepository,
	imageRepo repository.ImageRepository,
) *ProductServ {
	return &ProductServ{
		productRepo: productRepo,
		imageRepo:   imageRepo,
	}
}

func (s *ProductServ) Add(ctx context.Context, product *model.ProductWithImagesData) error {
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	prodID, err := s.productRepo.Add(ctx, tx, &product.Product)
	if err != nil {
		return err
	}

	var addedImages []struct {
		imgID     int64
		objectKey string
	}

	for _, image := range product.ProductImagesData {
		imgID, objectKey, err := s.imageRepo.Add(ctx, tx, prodID, &image)
		if err != nil {
			for _, added := range addedImages {
				_ = s.imageRepo.DeleteByID(ctx, added.imgID, added.objectKey)
			}
			return err
		}

		addedImages = append(addedImages, struct {
			imgID     int64
			objectKey string
		}{imgID: imgID, objectKey: objectKey})
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ProductServ) GetByIDs(ctx context.Context, IDs []int64) ([]model.ProductWithImages, error) {
	products, err := s.productRepo.GetByIDs(ctx, IDs)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return []model.ProductWithImages{}, nil
	}

	images, err := s.imageRepo.GetByIDs(ctx, IDs)
	if err != nil {
		return nil, err
	}

	imagesMap := make(map[int64][]model.Image)
	for _, img := range images {
		imagesMap[img.ProductID] = append(imagesMap[img.ProductID], img)
	}

	result := make([]model.ProductWithImages, 0, len(products))
	for _, product := range products {
		result = append(result, model.ProductWithImages{
			Product:       product,
			ProductImages: imagesMap[product.ID],
		})
	}

	return result, nil
}
