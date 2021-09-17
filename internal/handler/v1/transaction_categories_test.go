package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/service"
	mockService "github.com/lotostudio/financial-api/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_listTransactionCategories(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactionCategories)

	categories := []domain.TransactionCategory{
		{
			ID:    1,
			Title: "food",
			Type:  domain.Expense,
		},
	}

	setResponseBody := func(categories []domain.TransactionCategory) string {
		body, _ := json.Marshal(categories)

		return string(body)
	}

	tests := []struct {
		name                 string
		_type                string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name:  "ok",
			_type: "",
			mockBehaviour: func(s *mockService.MockTransactionCategories) {
				s.EXPECT().List(context.Background()).Return(categories, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(categories),
		},
		{
			name:  "ok",
			_type: "expense",
			mockBehaviour: func(s *mockService.MockTransactionCategories) {
				s.EXPECT().ListByType(context.Background(), domain.Expense).Return(categories, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(categories),
		},
		{
			name:  "invalid type",
			_type: "qwe",
			mockBehaviour: func(s *mockService.MockTransactionCategories) {
				s.EXPECT().ListByType(context.Background(), domain.TransactionType("qwe")).Return(categories, domain.ErrInvalidTransactionType)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"invalid type of transaction"}`,
		},
		{
			name:  "error",
			_type: "",
			mockBehaviour: func(s *mockService.MockTransactionCategories) {
				s.EXPECT().List(context.Background()).Return(categories, errors.New("general error"))
			},
			expectedCodeStatus:   500,
			expectedResponseBody: `{"message":"general error"}`,
		},

		{
			name:  "error",
			_type: "expense",
			mockBehaviour: func(s *mockService.MockTransactionCategories) {
				s.EXPECT().ListByType(context.Background(), domain.Expense).Return(categories, errors.New("general error"))
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

			tcService := mockService.NewMockTransactionCategories(c)
			tt.mockBehaviour(tcService)

			services := &service.Services{TransactionCategories: tcService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/transaction-categories", handler.listTransactionCategories)

			// Create Request
			w := httptest.NewRecorder()
			queryString := ""

			if tt._type != "" {
				queryString = "?type=" + tt._type
			}

			req := httptest.NewRequest("GET", "/transaction-categories"+queryString, bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
