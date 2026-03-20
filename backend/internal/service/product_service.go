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

func (s *ProductServ) GetByIDs(ctx context.Context, IDs []int) ([]model.ProductWithImages, error) {
	products, err := s.productRepo.GetByIDs(ctx, IDs)
	if err != nil {
		return nil, err
	}

	var prodsWithImages []model.ProductWithImages
	for _, product := range products {
		var prodWithImages model.ProductWithImages

		images, err := s.imageRepo.GetByID(ctx, product.ID)
		if err != nil {
			return nil, err
		}
		prodWithImages.Product = product
		prodWithImages.ProductImages = images

		prodsWithImages = append(prodsWithImages, prodWithImages)
	}

	return prodsWithImages, nil
}
