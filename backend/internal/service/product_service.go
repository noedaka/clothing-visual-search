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
	id, err := s.productRepo.Add(ctx, &product.Product)
	if err != nil {
		return err
	}

	for _, image := range product.ProductImagesData {
		if err = s.imageRepo.Add(ctx, id, &image); err != nil {
			return err
		}
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
