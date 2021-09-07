package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type AccountTypesService struct {
	repo repo.AccountTypes
}

func newAccountTypesService(repo repo.AccountTypes) *AccountTypesService {
	return &AccountTypesService{
		repo: repo,
	}
}

func (s *AccountTypesService) List(ctx context.Context) ([]domain.AccountType, error) {
	return s.repo.List(ctx)
}
