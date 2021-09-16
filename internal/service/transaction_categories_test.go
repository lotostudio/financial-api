package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockTransactionCategoriesService(t *testing.T) (*TransactionCategoryService, *mockRepo.MockTransactionCategories) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	tcRepo := mockRepo.NewMockTransactionCategories(mockCtl)

	s := newTransactionCategoriesService(tcRepo)

	return s, tcRepo
}

func TestTransactionCategoryService_List(t *testing.T) {
	s, ctRepo := mockTransactionCategoriesService(t)

	ctx := context.Background()

	ctRepo.EXPECT().List(ctx).Return([]domain.TransactionCategory{}, nil)

	categories, err := s.List(ctx)

	require.NoError(t, err)
	require.IsType(t, []domain.TransactionCategory{}, categories)
}

func TestTransactionCategoryService_ListByType(t *testing.T) {
	s, ctRepo := mockTransactionCategoriesService(t)

	ctx := context.Background()

	ctRepo.EXPECT().ListByType(ctx, domain.Transfer).Return([]domain.TransactionCategory{}, nil)

	categories, err := s.ListByType(ctx, domain.Transfer)

	require.NoError(t, err)
	require.IsType(t, []domain.TransactionCategory{}, categories)
}

func TestTransactionCategoryService_ListByTypeInvalidType(t *testing.T) {
	s, _ := mockTransactionCategoriesService(t)

	ctx := context.Background()

	_, err := s.ListByType(ctx, "qwe")

	require.ErrorIs(t, err, domain.ErrInvalidTransactionType)
}
