package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockTransactionTypesService(t *testing.T) (*TransactionTypesService, *mockRepo.MockTransactionTypes) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	ttRepo := mockRepo.NewMockTransactionTypes(mockCtl)

	s := newTransactionTypesService(ttRepo)

	return s, ttRepo
}

func TestTransactionTypesService_List(t *testing.T) {
	s, ttRepo := mockTransactionTypesService(t)

	ctx := context.Background()

	ttRepo.EXPECT().List(ctx).Return([]domain.TransactionType{}, nil)

	types, err := s.List(ctx)

	require.NoError(t, err)
	require.IsType(t, []domain.TransactionType{}, types)
}
