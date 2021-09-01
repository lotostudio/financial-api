package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

//go:generate mockgen -source=repo.go -destination=mocks/mock.go

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (int64, error)
	Get(ctx context.Context, id int64) (domain.User, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error)
}

type Currencies interface {
	List(ctx context.Context) ([]domain.Currency, error)
	Get(ctx context.Context, id int) (domain.Currency, error)
}

type Accounts interface {
	List(ctx context.Context, userID int64) ([]domain.Account, error)
	Create(ctx context.Context, toCreate domain.AccountToCreate, userID int64, currencyID int) (domain.Account, error)
	Get(ctx context.Context, id int64) (domain.Account, error)
	Update(ctx context.Context, toUpdate domain.AccountToUpdate, id int64, _type string) (domain.Account, error)
	Delete(ctx context.Context, id int64) error
}

type Repos struct {
	Users
	Currencies
	Accounts
}

func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{
		Users:      newUsersRepo(db),
		Currencies: newCurrenciesRepo(db),
		Accounts:   newAccountsRepo(db),
	}
}
