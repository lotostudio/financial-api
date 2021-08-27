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

type Repos struct {
	Users
}

func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{
		Users: newUsersRepo(db),
	}
}
