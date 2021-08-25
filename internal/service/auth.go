package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/pkg/auth"
	"github.com/lotostudio/financial-api/pkg/hash"
	"strconv"
)

type AuthService struct {
	repo         repo.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func newAuthService(repo repo.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *AuthService {
	return &AuthService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Register(ctx context.Context, toCreate domain.UserToCreate) (domain.User, error) {
	passwordHash, err := s.hasher.Hash(toCreate.Password)

	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		Email:     toCreate.Email,
		FirstName: toCreate.FirstName,
		LastName:  toCreate.LastName,
		Password:  passwordHash,
	}

	userId, err := s.repo.Create(ctx, user)

	if err != nil {
		return domain.User{}, err
	}

	return s.repo.Get(ctx, userId)
}

func (s *AuthService) Login(ctx context.Context, toLogin domain.UserToLogin) (domain.Tokens, error) {
	passwordHash, err := s.hasher.Hash(toLogin.Password)

	if err != nil {
		return domain.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, toLogin.Email, passwordHash)

	if err != nil {
		return domain.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *AuthService) createSession(ctx context.Context, userId int64) (domain.Tokens, error) {
	var res domain.Tokens
	var err error

	res.AccessToken, err = s.tokenManager.Issue(strconv.FormatInt(userId, 10))

	return res, err
}
