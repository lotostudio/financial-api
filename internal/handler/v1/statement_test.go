package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/service"
	mockService "github.com/lotostudio/financial-api/internal/service/mocks"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_getStatement(t *testing.T) {
	type mockBehaviour func(s *mockService.MockStats, a *mockService.MockAccounts)

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockStats, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, nil)
				s.EXPECT().Statement(context.Background(), gomock.Any()).Return(domain.Statement{}, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: `{"account":{"id":0,"title":"","balance":0,"currency":"","type":"","createdAt":"0001-01-01T00:00:00Z"},"balanceIn":{"date":"0001-01-01T00:00:00Z","value":0},"balanceOut":{"date":"0001-01-01T00:00:00Z","value":0},"transactions":null}`,
		},
		{
			name: "not found",
			mockBehaviour: func(s *mockService.MockStats, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, repo.ErrAccountNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account doesn't exists"}`,
		},
		{
			name: "forbidden",
			mockBehaviour: func(s *mockService.MockStats, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, service.ErrAccountForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"account forbidden to access"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockStats, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, errors.New("general error"))
			},
			expectedCodeStatus:   500,
			expectedResponseBody: `{"message":"general error"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockStats, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, nil)
				s.EXPECT().Statement(context.Background(), gomock.Any()).Return(domain.Statement{}, errors.New("general error"))
			},
			expectedCodeStatus:   500,
			expectedResponseBody: `{"message":"general error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			tService := mockService.NewMockStats(c)
			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(tService, aService)

			services := &service.Services{Stats: tService, Accounts: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/accounts/:id/statement", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.getStatement)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/accounts/%d/statement", accountID),
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
