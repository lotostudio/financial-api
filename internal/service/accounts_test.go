package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	userId    = int64(1)
	accountId = int64(2)
)

func mockAccountsService(t *testing.T) (*AccountsService, *mockRepo.MockAccounts, *mockRepo.MockCurrencies) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	aRepo := mockRepo.NewMockAccounts(mockCtl)
	cRepo := mockRepo.NewMockCurrencies(mockCtl)

	s := newAccountsService(aRepo, cRepo)

	return s, aRepo, cRepo
}

func TestAccountsService_List(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().List(ctx, userId).Return([]domain.Account{}, nil)

	accounts, err := s.List(ctx, userId)

	require.NoError(t, err)
	require.IsType(t, []domain.Account{}, accounts)
}

func TestAccountsService_Create(t *testing.T) {
	s, aRepo, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    "",
		Number:  nil,
		Term:    nil,
		Rate:    nil,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)
	aRepo.EXPECT().Create(ctx, toCreate, userId, 1).Return(domain.Account{}, nil)

	account, err := s.Create(ctx, toCreate, userId, 1)

	require.NoError(t, err)
	require.IsType(t, domain.Account{}, account)
}

func TestAccountsService_CreateCurrencyNotFound(t *testing.T) {
	s, _, cRepo := mockAccountsService(t)

	ctx := context.Background()

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{}, repo.ErrCurrencyNotFound)

	_, err := s.Create(ctx, domain.AccountToCreate{}, userId, 1)

	require.ErrorIs(t, err, repo.ErrCurrencyNotFound)
}

func TestAccountsService_CreateInvalidLoanData(t *testing.T) {
	s, _, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    "loan",
		Number:  nil,
		Term:    nil,
		Rate:    nil,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.ErrorIs(t, err, ErrInvalidLoanData)
}

func TestAccountsService_CreateInvalidDepositData(t *testing.T) {
	s, _, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    "deposit",
		Number:  nil,
		Term:    nil,
		Rate:    nil,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.ErrorIs(t, err, ErrInvalidDepositData)
}

func TestAccountsService_CreateInvalidCardData(t *testing.T) {
	s, _, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    "card",
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.ErrorIs(t, err, ErrInvalidCardData)
}

func TestAccountsService_CreateGeneralError(t *testing.T) {
	s, aRepo, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "qwe",
		Balance: 123,
		Type:    "cash",
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)
	aRepo.EXPECT().Create(ctx, toCreate, userId, 1).Return(domain.Account{}, errors.New("general error"))

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.Error(t, err)
}

func TestAccountsService_Get(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
	}, nil)

	account, err := s.Get(ctx, accountId, userId)

	require.NoError(t, err)
	require.IsType(t, domain.Account{}, account)
}

func TestAccountsService_GetAccountNotFound(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{}, repo.ErrAccountNotFound)

	_, err := s.Get(ctx, accountId, userId)

	require.ErrorIs(t, err, repo.ErrAccountNotFound)
}

func TestAccountsService_GetAccountForbidden(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: int64(99),
	}, nil)

	_, err := s.Get(ctx, accountId, userId)

	require.ErrorIs(t, err, ErrAccountForbidden)
}

func TestAccountsService_Update(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()
	title, balance := "title", 12.2
	toUpdate := domain.AccountToUpdate{
		Title:   &title,
		Balance: &balance,
		Number:  nil,
		Term:    nil,
		Rate:    nil,
	}

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    "cash",
	}, nil)
	aRepo.EXPECT().Update(ctx, toUpdate, accountId, "cash").Return(domain.Account{}, nil)

	account, err := s.Update(ctx, toUpdate, accountId, userId)

	require.NoError(t, err)
	require.IsType(t, domain.Account{}, account)
}

func TestAccountsService_UpdateInstanceError(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{}, repo.ErrAccountNotFound)

	_, err := s.Update(ctx, domain.AccountToUpdate{}, accountId, userId)

	require.ErrorIs(t, err, repo.ErrAccountNotFound)
}

func TestAccountsService_UpdateInvalidLoanData(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    "loan",
	}, nil)

	_, err := s.Update(ctx, domain.AccountToUpdate{}, accountId, userId)

	require.ErrorIs(t, err, ErrInvalidLoanData)
}

func TestAccountsService_UpdateInvalidDepositData(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    "deposit",
	}, nil)

	_, err := s.Update(ctx, domain.AccountToUpdate{}, accountId, userId)

	require.ErrorIs(t, err, ErrInvalidDepositData)
}

func TestAccountsService_UpdateInvalidCardData(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    "card",
	}, nil)

	_, err := s.Update(ctx, domain.AccountToUpdate{}, accountId, userId)

	require.ErrorIs(t, err, ErrInvalidCardData)
}

func TestAccountsService_UpdateGeneralError(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()
	title, balance := "title", 12.2
	toUpdate := domain.AccountToUpdate{
		Title:   &title,
		Balance: &balance,
		Number:  nil,
		Term:    nil,
		Rate:    nil,
	}

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    "cash",
	}, nil)
	aRepo.EXPECT().Update(ctx, toUpdate, accountId, "cash").Return(domain.Account{}, errors.New("general error"))

	_, err := s.Update(ctx, toUpdate, accountId, userId)

	require.Error(t, err)
}

func TestAccountsService_Delete(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
	}, nil)
	aRepo.EXPECT().Delete(ctx, accountId).Return(nil)

	err := s.Delete(ctx, accountId, userId)

	require.NoError(t, err)
}

func TestAccountsService_DeleteInstanceError(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{}, repo.ErrAccountNotFound)

	err := s.Delete(ctx, accountId, userId)

	require.Error(t, err, repo.ErrAccountNotFound)
}
