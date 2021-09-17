package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type TransactionTypesService struct {
	repo repo.TransactionTypes
}

func newTransactionTypesService(repo repo.TransactionTypes) *TransactionTypesService {
	return &TransactionTypesService{
		repo: repo,
	}
}

func (s *TransactionTypesService) List(ctx context.Context) ([]domain.TransactionType, error) {
	return s.repo.List(ctx)
}
