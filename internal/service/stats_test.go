package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func mockStatsService(t *testing.T) (Stats, *mockRepo.MockAccounts, *mockRepo.MockBalances, *mockRepo.MockTransactions) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	aRepo := mockRepo.NewMockAccounts(mockCtl)
	bRepo := mockRepo.NewMockBalances(mockCtl)
	tRepo := mockRepo.NewMockTransactions(mockCtl)

	s := newStatsService(aRepo, bRepo, tRepo)

	return s, aRepo, bRepo, tRepo
}

func TestStatsService_Statement(t *testing.T) {
	s, aRepo, bRepo, tRepo := mockStatsService(t)

	accId := int64(1)
	dateFrom, dateTo := time.Now(), time.Now().Add(24*time.Hour)

	ctx := context.Background()
	filter := domain.TransactionsFilter{
		AccountId:   &accId,
		CreatedFrom: &dateFrom,
		CreatedTo:   &dateTo,
	}

	aRepo.EXPECT().Get(gomock.Any(), *filter.AccountId).Return(domain.Account{}, nil)
	bRepo.EXPECT().Get(gomock.Any(), *filter.AccountId, *filter.CreatedFrom).Return(domain.Balance{}, nil)
	bRepo.EXPECT().Get(gomock.Any(), *filter.AccountId, *filter.CreatedTo).Return(domain.Balance{}, nil)
	tRepo.EXPECT().List(gomock.Any(), filter).Return([]domain.Transaction{}, nil)

	st, err := s.Statement(ctx, filter)

	require.NoError(t, err)
	require.IsType(t, domain.Statement{}, st)
}
