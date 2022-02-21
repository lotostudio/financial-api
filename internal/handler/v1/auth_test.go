package v1

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/service"
	mockService "github.com/lotostudio/financial-api/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_register(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuth, user domain.UserToCreate)

	tests := []struct {
		name          string
		requestBody   string
		requestUser   domain.UserToCreate
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:        "ok",
			requestBody: `{"firstName": "Sirius", "lastName": "Sam", "email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserToCreate{
				FirstName: "Sirius",
				LastName:  "Sam",
				Email:     "qweqweqwe@gmail.com",
				Password:  "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToCreate) {
				s.EXPECT().Register(context.Background(), user).Return(domain.User{
					ID:        1,
					Email:     "qweqweqwe@gmail.com",
					FirstName: "Sirius",
					LastName:  "Sam",
					Password:  "qweqweqwe",
				}, nil)
			},
			statusCode:   201,
			responseBody: `{"id":1,"email":"qweqweqwe@gmail.com","firstName":"Sirius","lastName":"Sam"}`,
		},
		{
			name:          "invalid request body",
			requestBody:   `{"firstName": "Sirius", "lastName": "Sam", "email": "qweqweqwe@gmail.com", "password": ""}`,
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToCreate) {},
			statusCode:    400,
			responseBody:  `{"message":"invalid request body - Key: 'UserToCreate.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:        "user already exists",
			requestBody: `{"firstName": "Sirius", "lastName": "Sam", "email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserToCreate{
				FirstName: "Sirius",
				LastName:  "Sam",
				Email:     "qweqweqwe@gmail.com",
				Password:  "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToCreate) {
				s.EXPECT().Register(context.Background(), user).Return(domain.User{}, repo.ErrUserAlreadyExists)
			},
			statusCode:   400,
			responseBody: `{"message":"user already exists"}`,
		},
		{
			name:        "error",
			requestBody: `{"firstName": "Sirius", "lastName": "Sam", "email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserToCreate{
				FirstName: "Sirius",
				LastName:  "Sam",
				Email:     "qweqweqwe@gmail.com",
				Password:  "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToCreate) {
				s.EXPECT().Register(context.Background(), user).Return(domain.User{}, errors.New("general error"))
			},
			statusCode:   500,
			responseBody: `{"message":"general error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuth(c)
			tt.mockBehaviour(auth, tt.requestUser)

			services := &service.Services{Auth: auth}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/register", handler.register)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_login(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuth, user domain.UserToLogin)

	tests := []struct {
		name          string
		requestBody   string
		requestUser   domain.UserToLogin
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:        "ok",
			requestBody: `{"email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserToLogin{
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToLogin) {
				s.EXPECT().Login(context.Background(), user).Return(domain.Tokens{
					AccessToken:  "token",
					RefreshToken: "token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"token","refreshToken":"token"}`,
		},
		{
			name:          "invalid request body",
			requestBody:   `{"email": "qweqweqwe", "password": "qweqweqwe"}`,
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToLogin) {},
			statusCode:    400,
			responseBody:  `{"message":"invalid request body - Key: 'UserToLogin.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{
			name:        "user does not exists",
			requestBody: `{"email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserToLogin{
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToLogin) {
				s.EXPECT().Login(context.Background(), user).Return(domain.Tokens{}, repo.ErrUserNotFound)
			},
			statusCode:   400,
			responseBody: `{"message":"user doesn't exists"}`,
		},
		{
			name:        "error",
			requestBody: `{"email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserToLogin{
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserToLogin) {
				s.EXPECT().Login(context.Background(), user).Return(domain.Tokens{}, errors.New("general error"))
			},
			statusCode:   500,
			responseBody: `{"message":"general error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuth(c)
			tt.mockBehaviour(auth, tt.requestUser)

			services := &service.Services{Auth: auth}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/login", handler.login)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_refresh(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuth)

	token := "token"

	tests := []struct {
		name          string
		authHeader    string
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:       "ok",
			authHeader: "Bearer " + token,
			mockBehaviour: func(s *mockService.MockAuth) {
				s.EXPECT().Refresh(context.Background(), token).Return(domain.Tokens{
					AccessToken:  "token",
					RefreshToken: "token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"token","refreshToken":"token"}`,
		},
		{
			name:       "expired token",
			authHeader: "Bearer " + token,
			mockBehaviour: func(s *mockService.MockAuth) {
				s.EXPECT().Refresh(context.Background(), token).Return(domain.Tokens{}, service.ErrRefreshTokenExpired)
			},
			statusCode:   400,
			responseBody: `{"message":"refresh token expired"}`,
		},
		{
			name:          "invalid token",
			authHeader:    "",
			mockBehaviour: func(s *mockService.MockAuth) {},
			statusCode:    401,
			responseBody:  `{"message":"empty auth header"}`,
		},
		{
			name:       "error",
			authHeader: "Bearer " + token,
			mockBehaviour: func(s *mockService.MockAuth) {
				s.EXPECT().Refresh(context.Background(), token).Return(domain.Tokens{}, errors.New("general error"))
			},
			statusCode:   500,
			responseBody: `{"message":"general error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuth(c)
			tt.mockBehaviour(auth)

			services := &service.Services{Auth: auth}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/refresh", handler.refresh)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(""))
			req.Header.Add("Authorization", tt.authHeader)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
