package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/pkg/hash"
)

type UsersService struct {
	repo   repo.Users
	hasher hash.PasswordHasher
}

func newUsersService(repo repo.Users, hasher hash.PasswordHasher) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *UsersService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UsersService) Get(ctx context.Context, userID int64) (domain.User, error) {
	return s.repo.Get(ctx, userID)
}

func (s *UsersService) UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error) {
	if toUpdate.Password != nil {
		passwordHash, err := s.hasher.Hash(*toUpdate.Password)

		if err != nil {
			return domain.User{}, err
		}

		toUpdate.Password = &passwordHash
	}

	return s.repo.UpdatePassword(ctx, userID, toUpdate)
}
