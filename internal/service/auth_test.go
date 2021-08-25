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

func mockAuthService(t *testing.T) (*AuthService, *mockRepo.MockUsers) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	usersRepo := mockRepo.NewMockUsers(mockCtl)
	authManager, _ := auth.NewJWTManager("key", time.Duration(1)*time.Hour)

	service := newAuthService(usersRepo, hash.NewSHA1PasswordHasher(""), authManager)

	return service, usersRepo
}

func TestAuthService_Register(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().Create(ctx, gomock.Any()).Return(int64(1), nil)
	usersRepo.EXPECT().Get(ctx, int64(1)).Return(domain.User{}, nil)

	user, err := service.Register(ctx, domain.UserToCreate{})

	require.NoError(t, err)
	require.IsType(t, domain.User{}, user)
}

func TestAuthService_Login(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, nil)

	res, err := service.Login(ctx, domain.UserToLogin{})

	require.NoError(t, err)
	require.IsType(t, domain.Tokens{}, res)
}

func TestAuthService_LoginErrUserNotExists(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, repo.ErrUserNotFound)

	_, err := service.Login(ctx, domain.UserToLogin{})

	require.ErrorIs(t, err, repo.ErrUserNotFound)
}

func TestAuthService_LoginErr(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, errDefault)

	_, err := service.Login(ctx, domain.UserToLogin{})

	require.Error(t, err)
}
