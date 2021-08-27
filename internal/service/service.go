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
	UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error)
}

type Auth interface {
	Register(ctx context.Context, user domain.UserToCreate) (domain.User, error)
	Login(ctx context.Context, user domain.UserToLogin) (domain.Tokens, error)
}

type Services struct {
	Users
	Auth
}

func NewServices(repos *repo.Repos, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *Services {
	return &Services{
		Users: newUsersService(repos.Users, hasher),
		Auth:  newAuthService(repos.Users, hasher, tokenManager),
	}
}
