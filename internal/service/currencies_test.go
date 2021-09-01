package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockCurrenciesService(t *testing.T) (*CurrenciesService, *mockRepo.MockCurrencies) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	cRepo := mockRepo.NewMockCurrencies(mockCtl)

	s := newCurrenciesService(cRepo)

	return s, cRepo
}

func TestCurrenciesService_List(t *testing.T) {
	s, cRepo := mockCurrenciesService(t)

	ctx := context.Background()

	cRepo.EXPECT().List(ctx).Return([]domain.Currency{}, nil)

	accounts, err := s.List(ctx)

	require.NoError(t, err)
	require.IsType(t, []domain.Currency{}, accounts)
}
