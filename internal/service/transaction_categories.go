package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type TransactionCategoryService struct {
	repo repo.TransactionCategories
}

func newTransactionCategoriesService(repo repo.TransactionCategories) *TransactionCategoryService {
	return &TransactionCategoryService{
		repo: repo,
	}
}

func (s *TransactionCategoryService) List(ctx context.Context) ([]domain.TransactionCategory, error) {
	return s.repo.List(ctx)
}

func (s *TransactionCategoryService) ListByType(ctx context.Context, _type domain.TransactionType) ([]domain.TransactionCategory, error) {
	if err := _type.Validate(); err != nil {
		return nil, err
	}

	return s.repo.ListByType(ctx, _type)
}
