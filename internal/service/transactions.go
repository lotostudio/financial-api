package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type TransactionsService struct {
	repo           repo.Transactions
	accountsRepo   repo.Accounts
	categoriesRepo repo.TransactionCategories
}

func newTransactionsService(repo repo.Transactions, accountsRepo repo.Accounts, categoriesRepo repo.TransactionCategories) *TransactionsService {
	return &TransactionsService{
		repo:           repo,
		accountsRepo:   accountsRepo,
		categoriesRepo: categoriesRepo,
	}
}

func (s *TransactionsService) List(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	return s.repo.List(ctx, userID)
}

func (s *TransactionsService) Create(ctx context.Context, toCreate domain.TransactionToCreate, userID int64,
	categoryId *int64, creditId *int64, debitId *int64) (domain.Transaction, error) {

	if err := toCreate.Type.Validate(); err != nil {
		return domain.Transaction{}, err
	}

	var category domain.TransactionCategory
	var err error

	if toCreate.Type != domain.Transfer {
		category, err = s.categoriesRepo.Get(ctx, *categoryId)

		if err != nil {
			return domain.Transaction{}, err
		}

		if category.Type != toCreate.Type {
			return domain.Transaction{}, ErrTransactionAndCategoryTypesMismatch
		}
	}

	switch toCreate.Type {
	case domain.Income:
		if _, err = s.checkIncome(ctx, userID, debitId); err != nil {
			return domain.Transaction{}, err
		}
	case domain.Expense:
		if _, err = s.checkExpense(ctx, userID, creditId); err != nil {
			return domain.Transaction{}, err
		}
	case domain.Transfer:
		if err = s.checkTransfer(ctx, userID, creditId, debitId); err != nil {
			return domain.Transaction{}, err
		}
	}

	transaction, err := s.repo.Create(ctx, toCreate, categoryId, creditId, debitId)

	if err != nil {
		return domain.Transaction{}, err
	}

	if category.Title != "" {
		transaction.Category = &category.Title
	}

	if transaction.Type == domain.Income || transaction.Type == domain.Transfer {
		debitAcc, err := s.accountsRepo.Get(ctx, *debitId)

		if err != nil {
			return transaction, err
		}

		transaction.Debit = &debitAcc
	}

	if transaction.Type == domain.Expense || transaction.Type == domain.Transfer {
		creditAcc, err := s.accountsRepo.Get(ctx, *creditId)

		if err != nil {
			return transaction, err
		}

		transaction.Credit = &creditAcc
	}

	return transaction, nil
}

func (s *TransactionsService) checkIncome(ctx context.Context, userID int64, debitId *int64) (domain.Account, error) {
	if debitId == nil {
		return domain.Account{}, ErrNoAccountSelected
	}

	debitAcc, err := s.accountsRepo.Get(ctx, *debitId)

	if err != nil {
		return debitAcc, err
	}

	if debitAcc.OwnerId != userID {
		return debitAcc, ErrDebitAccountForbidden
	}

	return debitAcc, nil
}

func (s *TransactionsService) checkExpense(ctx context.Context, userID int64, creditId *int64) (domain.Account, error) {
	if creditId == nil {
		return domain.Account{}, ErrNoAccountSelected
	}

	creditAcc, err := s.accountsRepo.Get(ctx, *creditId)

	if err != nil {
		return creditAcc, err
	}

	if creditAcc.OwnerId != userID {
		return creditAcc, ErrDebitAccountForbidden
	}

	return creditAcc, nil
}

func (s *TransactionsService) checkTransfer(ctx context.Context, userID int64, creditId *int64, debitId *int64) error {
	creditAcc, err := s.checkExpense(ctx, userID, creditId)

	if err != nil {
		return err
	}

	debitAcc, err := s.checkIncome(ctx, userID, debitId)

	if err != nil {
		return err
	}

	if creditAcc.Currency != debitAcc.Currency {
		return ErrAccountsHaveDifferenceCurrencies
	}

	return nil
}
