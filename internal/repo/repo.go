package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
	"time"
)

//go:generate mockgen -source=repo.go -destination=mocks/mock.go

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (int64, error)
	Get(ctx context.Context, id int64) (domain.User, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error)
}

type Sessions interface {
	Create(ctx context.Context, userID int64) error
	GetByToken(ctx context.Context, token string) (domain.Session, error)
	Update(ctx context.Context, toUpdate domain.SessionToUpdate, userID int64) (domain.Session, error)
}

type Currencies interface {
	List(ctx context.Context) ([]domain.Currency, error)
	Get(ctx context.Context, id int) (domain.Currency, error)
}

type Accounts interface {
	List(ctx context.Context, userID int64) ([]domain.Account, error)
	CountByTypes(ctx context.Context, userID int64, _type domain.AccountType, types ...domain.AccountType) (int64, error)
	Create(ctx context.Context, toCreate domain.AccountToCreate, userID int64, currencyID int) (domain.Account, error)
	Get(ctx context.Context, id int64) (domain.Account, error)
	Update(ctx context.Context, toUpdate domain.AccountToUpdate, id int64, _type domain.AccountType) (domain.Account, error)
	Delete(ctx context.Context, id int64) error
}

type AccountTypes interface {
	List(ctx context.Context) ([]domain.AccountType, error)
}

type Transactions interface {
	List(ctx context.Context, filter domain.TransactionsFilter) ([]domain.Transaction, error)
	Create(ctx context.Context, toCreate domain.TransactionToCreate, categoryId *int64, creditId *int64,
		debitId *int64) (domain.Transaction, error)
	GetOwner(ctx context.Context, id int64) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type TransactionCategories interface {
	List(ctx context.Context) ([]domain.TransactionCategory, error)
	ListByType(ctx context.Context, _type domain.TransactionType) ([]domain.TransactionCategory, error)
	Get(ctx context.Context, id int64) (domain.TransactionCategory, error)
}

type TransactionTypes interface {
	List(ctx context.Context) ([]domain.TransactionType, error)
}

type Balances interface {
	Get(ctx context.Context, accountID int64, date time.Time) (domain.Balance, error)
}

type Repos struct {
	Users
	Sessions
	Currencies
	Accounts
	AccountTypes
	Transactions
	TransactionCategories
	TransactionTypes
	Balances
}

func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{
		Users:                 newUsersRepo(db),
		Sessions:              newSessionsRepo(db),
		Currencies:            newCurrenciesRepo(db),
		Accounts:              newAccountsRepo(db),
		AccountTypes:          newAccountTypesRepo(db),
		Transactions:          newTransactionsRepo(db),
		TransactionCategories: newTransactionCategoriesRepo(db),
		TransactionTypes:      newTransactionTypesRepo(db),
		Balances:              newBalancesRepo(db),
	}
}
