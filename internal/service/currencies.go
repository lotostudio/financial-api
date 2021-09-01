package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type CurrenciesService struct {
	repo repo.Currencies
}

func newCurrenciesService(repo repo.Currencies) *CurrenciesService {
	return &CurrenciesService{
		repo: repo,
	}
}

func (s *CurrenciesService) List(ctx context.Context) ([]domain.Currency, error) {
	return s.repo.List(ctx)
}
