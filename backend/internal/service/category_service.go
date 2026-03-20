package service

import (
	"context"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
	"github.com/noedaka/clothing-visual-search/backend/internal/repository"
)

type CategoryServ struct {
	repo repository.CategoryRepository
}

func NewCategoryServ(repo repository.CategoryRepository) *CategoryServ {
	return &CategoryServ{repo: repo}
}

func (s *CategoryServ) Add(ctx context.Context, category string) error {
	return s.repo.Add(ctx, category)
}

func (s *CategoryServ) List(ctx context.Context) ([]model.Category, error) {
	return s.repo.List(ctx)
}
