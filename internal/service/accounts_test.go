package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/config"
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
	aCfg := config.Account{
		CardAndCashLimit:    1,
		LoanAndDepositLimit: 1,
	}

	s := newAccountsService(aRepo, cRepo, aCfg)

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

func TestAccountsService_ListGrouped(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()
	number := "0000"

	aRepo.EXPECT().List(ctx, userId).Return([]domain.Account{
		{
			ID:       1,
			Title:    "acc1",
			Balance:  12.1,
			Currency: "KZT",
			Type:     domain.Card,
			Number:   &number,
		},
		{
			ID:       2,
			Title:    "acc2",
			Balance:  12.1,
			Currency: "KZT",
			Type:     domain.Card,
			Number:   &number,
		},
	}, nil)

	accounts, err := s.ListGrouped(ctx, userId)

	require.NoError(t, err)
	require.IsType(t, domain.GroupedAccounts{}, accounts)
}

func TestAccountsService_ListGroupedError(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().List(ctx, userId).Return(nil, errors.New("general error"))

	_, err := s.ListGrouped(ctx, userId)

	require.Error(t, err)
}

func TestAccountsService_Create(t *testing.T) {
	s, aRepo, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    domain.Cash,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)
	aRepo.EXPECT().CountByTypes(ctx, userId, domain.Cash, domain.Card).Return(int64(0), nil)
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
		Type:    domain.Loan,
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
		Type:    domain.Deposit,
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
		Type:    domain.Card,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.ErrorIs(t, err, ErrInvalidCardData)
}

func TestAccountsService_CreateCashCardLimit(t *testing.T) {
	s, aRepo, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    domain.Cash,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)
	aRepo.EXPECT().CountByTypes(ctx, userId, domain.Cash, domain.Card).Return(int64(1), nil)

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.ErrorIs(t, err, ErrAccountCountLimited)
}

func TestAccountsService_CreateLoanDepositLimit(t *testing.T) {
	s, aRepo, cRepo := mockAccountsService(t)

	ctx := context.Background()
	var term uint8 = 1
	var rate float32 = 1.1
	toCreate := domain.AccountToCreate{
		Title:   "",
		Balance: 0,
		Type:    domain.Loan,
		Term:    &term,
		Rate:    &rate,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)
	aRepo.EXPECT().CountByTypes(ctx, userId, domain.Loan, domain.Deposit).Return(int64(1), nil)

	_, err := s.Create(ctx, toCreate, userId, 1)

	require.ErrorIs(t, err, ErrAccountCountLimited)
}

func TestAccountsService_CreateGeneralError(t *testing.T) {
	s, aRepo, cRepo := mockAccountsService(t)

	ctx := context.Background()
	toCreate := domain.AccountToCreate{
		Title:   "qwe",
		Balance: 123,
		Type:    domain.Cash,
	}

	cRepo.EXPECT().Get(ctx, 1).Return(domain.Currency{ID: 1, Code: "KZT"}, nil)
	aRepo.EXPECT().CountByTypes(ctx, userId, domain.Cash, domain.Card).Return(int64(0), nil)
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
	}

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    domain.Cash,
	}, nil)
	aRepo.EXPECT().Update(ctx, toUpdate, accountId, domain.Cash).Return(domain.Account{}, nil)

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
		Type:    domain.Loan,
	}, nil)

	_, err := s.Update(ctx, domain.AccountToUpdate{}, accountId, userId)

	require.ErrorIs(t, err, ErrInvalidLoanData)
}

func TestAccountsService_UpdateInvalidDepositData(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    domain.Deposit,
	}, nil)

	_, err := s.Update(ctx, domain.AccountToUpdate{}, accountId, userId)

	require.ErrorIs(t, err, ErrInvalidDepositData)
}

func TestAccountsService_UpdateInvalidCardData(t *testing.T) {
	s, aRepo, _ := mockAccountsService(t)

	ctx := context.Background()

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    domain.Card,
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
	}

	aRepo.EXPECT().Get(ctx, accountId).Return(domain.Account{
		OwnerId: userId,
		Type:    domain.Cash,
	}, nil)
	aRepo.EXPECT().Update(ctx, toUpdate, accountId, domain.Cash).Return(domain.Account{}, errors.New("general error"))

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
