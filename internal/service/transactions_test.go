package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockTransactionsService(t *testing.T) (*TransactionsService, *mockRepo.MockTransactions, *mockRepo.MockAccounts,
	*mockRepo.MockTransactionCategories) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	tRepo := mockRepo.NewMockTransactions(mockCtl)
	aRepo := mockRepo.NewMockAccounts(mockCtl)
	tcRepo := mockRepo.NewMockTransactionCategories(mockCtl)

	s := newTransactionsService(tRepo, aRepo, tcRepo)

	return s, tRepo, aRepo, tcRepo
}

func TestTransactionsService_List(t *testing.T) {
	s, tRepo, _, _ := mockTransactionsService(t)

	ctx := context.Background()
	filter := domain.TransactionsFilter{}

	tRepo.EXPECT().List(ctx, filter).Return([]domain.Transaction{}, nil)

	categories, err := s.List(ctx, filter)

	require.NoError(t, err)
	require.IsType(t, []domain.Transaction{}, categories)
}

func TestTransactionsService_CreateIncome(t *testing.T) {
	s, tRepo, aRepo, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Income,
	}
	var categoryId, debitId = new(int64), new(int64)
	*categoryId = 1
	*debitId = 1

	debit := domain.Account{
		OwnerId: userId,
	}

	tr := domain.Transaction{
		Type: domain.Income,
	}

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Income,
		Title: "salary",
	}, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)
	tRepo.EXPECT().Create(ctx, toCreate, categoryId, nil, debitId).Return(tr, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)

	created, err := s.Create(ctx, toCreate, userId, categoryId, nil, debitId)

	require.NoError(t, err)
	require.IsType(t, domain.Transaction{}, created)
}

func TestTransactionsService_CreateIncomeErrAccountNotSelected(t *testing.T) {
	s, _, _, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Income,
	}
	var categoryId = new(int64)
	*categoryId = 1

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Income,
		Title: "salary",
	}, nil)

	_, err := s.Create(ctx, toCreate, userId, categoryId, nil, nil)

	require.ErrorIs(t, err, ErrNoAccountSelected)
}

func TestTransactionsService_CreateIncomeErrAccountForbidden(t *testing.T) {
	s, _, aRepo, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Income,
	}
	var categoryId, debitId = new(int64), new(int64)
	*categoryId = 1
	*debitId = 1

	debit := domain.Account{
		OwnerId: userId + 1,
	}

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Income,
		Title: "salary",
	}, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)

	_, err := s.Create(ctx, toCreate, userId, categoryId, nil, debitId)

	require.ErrorIs(t, err, ErrDebitAccountForbidden)
}

func TestTransactionsService_CreateExpense(t *testing.T) {
	s, tRepo, aRepo, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Expense,
	}
	var categoryId, creditId = new(int64), new(int64)
	*categoryId = 1
	*creditId = 1

	credit := domain.Account{
		OwnerId: userId,
	}

	tr := domain.Transaction{
		Type: domain.Expense,
	}

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Expense,
		Title: "food",
	}, nil)
	aRepo.EXPECT().Get(ctx, *creditId).Return(credit, nil)
	tRepo.EXPECT().Create(ctx, toCreate, categoryId, creditId, nil).Return(tr, nil)
	aRepo.EXPECT().Get(ctx, *creditId).Return(credit, nil)

	created, err := s.Create(ctx, toCreate, userId, categoryId, creditId, nil)

	require.NoError(t, err)
	require.IsType(t, domain.Transaction{}, created)
}

func TestTransactionsService_CreateExpenseErrAccountNotSelected(t *testing.T) {
	s, _, _, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Expense,
	}
	var categoryId = new(int64)
	*categoryId = 1

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Expense,
		Title: "food",
	}, nil)

	_, err := s.Create(ctx, toCreate, userId, categoryId, nil, nil)

	require.ErrorIs(t, err, ErrNoAccountSelected)
}

func TestTransactionsService_CreateExpenseErrAccountForbidden(t *testing.T) {
	s, _, aRepo, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Expense,
	}
	var categoryId, creditId = new(int64), new(int64)
	*categoryId = 1
	*creditId = 1

	credit := domain.Account{
		OwnerId: userId + 1,
	}

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Expense,
		Title: "food",
	}, nil)
	aRepo.EXPECT().Get(ctx, *creditId).Return(credit, nil)

	_, err := s.Create(ctx, toCreate, userId, categoryId, creditId, nil)

	require.ErrorIs(t, err, ErrCreditAccountForbidden)
}

func TestTransactionsService_CreateTransfer(t *testing.T) {
	s, tRepo, aRepo, _ := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Transfer,
	}
	var categoryId, creditId, debitId = new(int64), new(int64), new(int64)
	*categoryId = 1
	*creditId = 1
	*debitId = 2

	credit := domain.Account{
		OwnerId: userId,
	}

	debit := domain.Account{
		OwnerId: userId,
	}

	tr := domain.Transaction{
		Type: domain.Transfer,
	}

	aRepo.EXPECT().Get(ctx, *creditId).Return(credit, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)
	tRepo.EXPECT().Create(ctx, toCreate, categoryId, creditId, debitId).Return(tr, nil)
	aRepo.EXPECT().Get(ctx, *creditId).Return(credit, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)

	created, err := s.Create(ctx, toCreate, userId, categoryId, creditId, debitId)

	require.NoError(t, err)
	require.IsType(t, domain.Transaction{}, created)
}

func TestTransactionsService_CreateTransferErrCurrenciesMismatch(t *testing.T) {
	s, _, aRepo, _ := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Transfer,
	}
	var categoryId, creditId, debitId = new(int64), new(int64), new(int64)
	*categoryId = 1
	*creditId = 1
	*debitId = 2

	credit := domain.Account{
		OwnerId:  userId,
		Currency: "KZT",
	}

	debit := domain.Account{
		OwnerId:  userId,
		Currency: "USD",
	}

	aRepo.EXPECT().Get(ctx, *creditId).Return(credit, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)

	_, err := s.Create(ctx, toCreate, userId, categoryId, creditId, debitId)

	require.ErrorIs(t, err, ErrAccountsHaveDifferenceCurrencies)
}

func TestTransactionsService_CreateInvalidType(t *testing.T) {
	s, _, _, _ := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: "qwe",
	}

	_, err := s.Create(ctx, toCreate, userId, nil, nil, nil)

	require.Error(t, err, domain.ErrInvalidTransactionType)
}

func TestTransactionsService_CreateErrType(t *testing.T) {
	s, _, _, ctRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Income,
	}

	var categoryId = new(int64)
	*categoryId = 1

	ctRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{}, errDefault)

	_, err := s.Create(ctx, toCreate, userId, categoryId, nil, nil)

	require.Error(t, err, errDefault)
}

func TestTransactionsService_CreateErrTypesMismatch(t *testing.T) {
	s, _, _, ctRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Income,
	}

	var categoryId = new(int64)
	*categoryId = 1

	ctRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type: domain.Expense,
	}, nil)

	_, err := s.Create(ctx, toCreate, userId, categoryId, nil, nil)

	require.Error(t, err, ErrTransactionAndCategoryTypesMismatch)
}

func TestTransactionsService_CreateErrDefault(t *testing.T) {
	s, tRepo, aRepo, tcRepo := mockTransactionsService(t)

	ctx := context.Background()
	toCreate := domain.TransactionToCreate{
		Type: domain.Income,
	}
	var categoryId, debitId = new(int64), new(int64)
	*categoryId = 1
	*debitId = 1

	debit := domain.Account{
		OwnerId: userId,
	}

	tr := domain.Transaction{
		Type: domain.Income,
	}

	tcRepo.EXPECT().Get(ctx, *categoryId).Return(domain.TransactionCategory{
		Type:  domain.Income,
		Title: "salary",
	}, nil)
	aRepo.EXPECT().Get(ctx, *debitId).Return(debit, nil)
	tRepo.EXPECT().Create(ctx, toCreate, categoryId, nil, debitId).Return(tr, errDefault)

	_, err := s.Create(ctx, toCreate, userId, categoryId, nil, debitId)

	require.ErrorIs(t, err, errDefault)
}

func TestTransactionsService_Delete(t *testing.T) {
	s, tRepo, _, _ := mockTransactionsService(t)

	ctx := context.Background()
	id := int64(1)

	tRepo.EXPECT().GetOwner(ctx, id).Return(userId, nil)
	tRepo.EXPECT().Delete(ctx, id).Return(nil)

	err := s.Delete(ctx, id, userId)

	require.NoError(t, err)
}

func TestTransactionsService_DeleteErrOwner(t *testing.T) {
	s, tRepo, _, _ := mockTransactionsService(t)

	ctx := context.Background()
	id := int64(1)

	tRepo.EXPECT().GetOwner(ctx, id).Return(userId, errDefault)

	err := s.Delete(ctx, id, userId)

	require.ErrorIs(t, err, errDefault)
}

func TestTransactionsService_DeleteErrForbidden(t *testing.T) {
	s, tRepo, _, _ := mockTransactionsService(t)

	ctx := context.Background()
	id := int64(1)

	tRepo.EXPECT().GetOwner(ctx, id).Return(userId+1, nil)

	err := s.Delete(ctx, id, userId)

	require.ErrorIs(t, err, ErrTransactionForbidden)
}

func TestTransactionsService_DeleteErr(t *testing.T) {
	s, tRepo, _, _ := mockTransactionsService(t)

	ctx := context.Background()
	id := int64(1)

	tRepo.EXPECT().GetOwner(ctx, id).Return(userId, nil)
	tRepo.EXPECT().Delete(ctx, id).Return(errDefault)

	err := s.Delete(ctx, id, userId)

	require.ErrorIs(t, err, errDefault)
}
