package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	mockRepo "github.com/lotostudio/financial-api/internal/repo/mocks"
	"github.com/lotostudio/financial-api/pkg/hash"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockUsersService(t *testing.T) (*UsersService, *mockRepo.MockUsers) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	usersRepo := mockRepo.NewMockUsers(mockCtl)

	service := newUsersService(usersRepo, hash.NewSHA1PasswordHasher(""))

	return service, usersRepo
}

func TestUsersService_List(t *testing.T) {
	service, usersRepo := mockUsersService(t)

	ctx := context.Background()

	usersRepo.EXPECT().List(ctx).Return([]domain.User{}, nil)

	res, err := service.List(ctx)

	require.NoError(t, err)
	require.IsType(t, []domain.User{}, res)
}

func TestUsersService_UpdatePassword(t *testing.T) {
	service, usersRepo := mockUsersService(t)

	ctx := context.Background()

	firstName, lastName, password := "Azamat", "Yergali", "qweqweqwe"
	hashedPassword, _ := service.hasher.Hash(password)

	usersRepo.EXPECT().UpdatePassword(ctx, int64(1), domain.UserToUpdate{
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &hashedPassword,
	}).Return(domain.User{}, nil)

	res, err := service.UpdatePassword(ctx, int64(1), domain.UserToUpdate{
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
	})

	require.NoError(t, err)
	require.IsType(t, domain.User{}, res)
}
