package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
)

type UsersService struct {
	repo repo.Users
}

func newUsersService(repo repo.Users) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}
