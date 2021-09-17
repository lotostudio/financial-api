package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/lotostudio/financial-api/pkg/auth"
	"github.com/lotostudio/financial-api/pkg/hash"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var errDefault = errors.New("error")

func mockAuthService(t *testing.T) (*AuthService, *mockRepo.MockUsers, *mockRepo.MockSessions) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	usersRepo := mockRepo.NewMockUsers(mockCtl)
	sRepo := mockRepo.NewMockSessions(mockCtl)
	authManager, _ := auth.NewJWTManager("key", time.Duration(1)*time.Hour, 32)

	service := newAuthService(usersRepo, sRepo, hash.NewSHA1PasswordHasher(""), authManager, 1*time.Second, 1*time.Second)

	return service, usersRepo, sRepo
}

func TestAuthService_Register(t *testing.T) {
	s, uRepo, sRepo := mockAuthService(t)

	ctx := context.Background()

	uRepo.EXPECT().Create(ctx, gomock.Any()).Return(userId, nil)
	uRepo.EXPECT().Get(ctx, userId).Return(domain.User{}, nil)
	sRepo.EXPECT().Create(ctx, userId)

	user, err := s.Register(ctx, domain.UserToCreate{})

	require.NoError(t, err)
	require.IsType(t, domain.User{}, user)
}

func TestAuthService_Login(t *testing.T) {
	s, uRepo, sRepo := mockAuthService(t)

	ctx := context.Background()

	uRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{
		ID: userId,
	}, nil)
	sRepo.EXPECT().Update(ctx, gomock.Any(), userId).Return(domain.Session{}, nil)

	res, err := s.Login(ctx, domain.UserToLogin{})

	require.NoError(t, err)
	require.IsType(t, domain.Tokens{}, res)
}

func TestAuthService_LoginErrUserNotExists(t *testing.T) {
	s, uRepo, _ := mockAuthService(t)

	ctx := context.Background()

	uRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, repo.ErrUserNotFound)

	_, err := s.Login(ctx, domain.UserToLogin{})

	require.ErrorIs(t, err, repo.ErrUserNotFound)
}

func TestAuthService_LoginErr(t *testing.T) {
	s, uRepo, _ := mockAuthService(t)

	ctx := context.Background()

	uRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, errDefault)

	_, err := s.Login(ctx, domain.UserToLogin{})

	require.ErrorIs(t, err, errDefault)
}

func TestAuthService_Refresh(t *testing.T) {
	s, _, sRepo := mockAuthService(t)

	ctx := context.Background()

	sRepo.EXPECT().GetByToken(ctx, "token").Return(domain.Session{
		ExpiresAt: time.Now().Add(1 * time.Hour),
		UserId:    userId,
	}, nil)
	sRepo.EXPECT().Update(ctx, gomock.Any(), userId).Return(domain.Session{}, nil)

	tokens, err := s.Refresh(ctx, "token")

	require.NoError(t, err)
	require.IsType(t, domain.Tokens{}, tokens)
}

func TestAuthService_RefreshExpiredToken(t *testing.T) {
	s, _, sRepo := mockAuthService(t)

	ctx := context.Background()

	sRepo.EXPECT().GetByToken(ctx, "token").Return(domain.Session{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		UserId:    userId,
	}, nil)

	_, err := s.Refresh(ctx, "token")

	require.ErrorIs(t, err, ErrRefreshTokenExpired)
}
