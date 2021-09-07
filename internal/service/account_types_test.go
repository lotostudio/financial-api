package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockAccountTypesService(t *testing.T) (*AccountTypesService, *mockRepo.MockAccountTypes) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	tRepo := mockRepo.NewMockAccountTypes(mockCtl)

	s := newAccountTypesService(tRepo)

	return s, tRepo
}

func TestAccountTypesService_List(t *testing.T) {
	s, tRepo := mockAccountTypesService(t)

	ctx := context.Background()

	tRepo.EXPECT().List(ctx).Return([]domain.AccountType{}, nil)

	types, err := s.List(ctx)

	require.NoError(t, err)
	require.IsType(t, []domain.AccountType{}, types)
}
