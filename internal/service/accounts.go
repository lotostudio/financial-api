package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type AccountsService struct {
	repo           repo.Accounts
	currenciesRepo repo.Currencies
}

func newAccountsService(repo repo.Accounts, currenciesRepo repo.Currencies) *AccountsService {
	return &AccountsService{
		repo:           repo,
		currenciesRepo: currenciesRepo,
	}
}

func (s *AccountsService) List(ctx context.Context, userID int64) ([]domain.Account, error) {
	return s.repo.List(ctx, userID)
}

func (s *AccountsService) ListGrouped(ctx context.Context, userID int64) (domain.GroupedAccounts, error) {
	accounts, err := s.repo.List(ctx, userID)

	if err != nil {
		return domain.GroupedAccounts{}, err
	}

	grouped := make(domain.GroupedAccounts)

	// Iterate through accounts and group them by types
	for _, a := range accounts {
		if val, ok := grouped[a.Type]; ok {
			grouped[a.Type] = append(val, a)
			continue
		}

		grouped[a.Type] = []domain.Account{a}
	}

	return grouped, nil
}

func (s *AccountsService) Create(ctx context.Context, toCreate domain.AccountToCreate, userID int64, currencyID int) (domain.Account, error) {
	currency, err := s.currenciesRepo.Get(ctx, currencyID)

	if err != nil {
		return domain.Account{}, err
	}

	// Check for required fields of loan account
	if toCreate.Type == domain.Loan && (toCreate.Term == nil || toCreate.Rate == nil) {
		return domain.Account{}, ErrInvalidLoanData
	}

	// Check for required fields of deposit account
	if toCreate.Type == domain.Deposit && (toCreate.Term == nil || toCreate.Rate == nil) {
		return domain.Account{}, ErrInvalidDepositData
	}

	// Check for required fields of card account
	if toCreate.Type == domain.Card && toCreate.Number == nil {
		return domain.Account{}, ErrInvalidCardData
	}

	account, err := s.repo.Create(ctx, toCreate, userID, currencyID)

	if err != nil {
		return account, err
	}

	account.Currency = currency.Code

	return account, nil
}

func (s *AccountsService) Get(ctx context.Context, id int64, userID int64) (domain.Account, error) {
	account, err := s.repo.Get(ctx, id)

	if err != nil {
		return account, err
	}

	if account.OwnerId != userID {
		return account, ErrAccountForbidden
	}

	return account, nil
}

func (s *AccountsService) Update(ctx context.Context, toUpdate domain.AccountToUpdate, id int64, userID int64) (domain.Account, error) {
	instance, err := s.Get(ctx, id, userID)

	if err != nil {
		return instance, err
	}

	// Check for required fields of loan account
	if instance.Type == domain.Loan && (toUpdate.Term == nil || toUpdate.Rate == nil) {
		return instance, ErrInvalidLoanData
	}

	// Check for required fields of deposit account
	if instance.Type == domain.Deposit && (toUpdate.Term == nil || toUpdate.Rate == nil) {
		return instance, ErrInvalidDepositData
	}

	// Check for required fields of card account
	if instance.Type == domain.Card && toUpdate.Number == nil {
		return domain.Account{}, ErrInvalidCardData
	}

	account, err := s.repo.Update(ctx, toUpdate, id, instance.Type)

	if err != nil {
		return account, err
	}

	account.Currency = instance.Currency

	return account, nil
}

func (s *AccountsService) Delete(ctx context.Context, id int64, userID int64) error {
	_, err := s.Get(ctx, id, userID)

	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
