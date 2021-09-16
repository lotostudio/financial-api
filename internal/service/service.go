package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/pkg/auth"
	"github.com/lotostudio/financial-api/pkg/hash"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
	Get(ctx context.Context, userID int64) (domain.User, error)
	UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error)
}

type Auth interface {
	Register(ctx context.Context, user domain.UserToCreate) (domain.User, error)
	Login(ctx context.Context, user domain.UserToLogin) (domain.Tokens, error)
}

type Currencies interface {
	List(ctx context.Context) ([]domain.Currency, error)
}

type Accounts interface {
	List(ctx context.Context, userID int64) ([]domain.Account, error)
	ListGrouped(ctx context.Context, userID int64) (domain.GroupedAccounts, error)
	Create(ctx context.Context, toCreate domain.AccountToCreate, userID int64, currencyID int) (domain.Account, error)
	Get(ctx context.Context, id int64, userID int64) (domain.Account, error)
	Update(ctx context.Context, toUpdate domain.AccountToUpdate, id int64, userID int64) (domain.Account, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type AccountTypes interface {
	List(ctx context.Context) ([]domain.AccountType, error)
}

type Transactions interface {
	List(ctx context.Context, filter domain.TransactionsFilter) ([]domain.Transaction, error)
	Create(ctx context.Context, toCreate domain.TransactionToCreate, userID int64, categoryId *int64, creditId *int64,
		debitId *int64) (domain.Transaction, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type TransactionCategories interface {
	List(ctx context.Context) ([]domain.TransactionCategory, error)
	ListByType(ctx context.Context, _type domain.TransactionType) ([]domain.TransactionCategory, error)
}

type Services struct {
	Users
	Auth
	Currencies
	Accounts
	AccountTypes
	Transactions
	TransactionCategories
}

func NewServices(repos *repo.Repos, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *Services {
	return &Services{
		Users:                 newUsersService(repos.Users, hasher),
		Auth:                  newAuthService(repos.Users, hasher, tokenManager),
		Currencies:            newCurrenciesService(repos.Currencies),
		Accounts:              newAccountsService(repos.Accounts, repos.Currencies),
		AccountTypes:          newAccountTypesService(repos.AccountTypes),
		Transactions:          newTransactionsService(repos.Transactions, repos.Accounts, repos.TransactionCategories),
		TransactionCategories: newTransactionCategoriesService(repos.TransactionCategories),
	}
}
