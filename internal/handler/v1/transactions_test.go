package v1

import (
	"bytes"
	"context"
	"encoding/json"
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
	"time"
)

const (
	transactionID = int64(5)
)

func TestHandler_listTransactions(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactions)

	transactions := []domain.Transaction{
		{
			ID:        1,
			Amount:    12.1,
			Type:      domain.Income,
			CreatedAt: time.Now(),
		},
	}

	setResponseBody := func(transactions []domain.Transaction) string {
		body, _ := json.Marshal(transactions)

		return string(body)
	}

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().List(context.Background(), gomock.Any()).Return(transactions, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(transactions),
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().List(context.Background(), gomock.Any()).Return(transactions, errors.New("general error"))
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

			tService := mockService.NewMockTransactions(c)
			tt.mockBehaviour(tService)

			services := &service.Services{Transactions: tService}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/transactions", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.listTransactions)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/transactions", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_transactionStats(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactions)

	stats := []domain.TransactionStat{
		{
			Category: "food",
			Value:    123,
		},
		{
			Category: "family",
			Value:    1000,
		},
	}

	setResponseBody := func(stats []domain.TransactionStat) string {
		body, _ := json.Marshal(stats)

		return string(body)
	}

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Stats(context.Background(), gomock.Any()).Return(stats, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(stats),
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Stats(context.Background(), gomock.Any()).Return(stats, errors.New("general error"))
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

			tService := mockService.NewMockTransactions(c)
			tt.mockBehaviour(tService)

			services := &service.Services{Transactions: tService}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/transactions/stats", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.transactionStats)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/transactions/stats", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_listTransactionsOfAccount(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactions, a *mockService.MockAccounts)

	transactions := []domain.Transaction{
		{
			ID:        1,
			Amount:    12.1,
			Type:      domain.Income,
			CreatedAt: time.Now(),
		},
	}

	setResponseBody := func(transactions []domain.Transaction) string {
		body, _ := json.Marshal(transactions)

		return string(body)
	}

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockTransactions, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, nil)
				s.EXPECT().List(context.Background(), gomock.Any()).Return(transactions, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(transactions),
		},
		{
			name: "not found",
			mockBehaviour: func(s *mockService.MockTransactions, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, repo.ErrAccountNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account doesn't exists"}`,
		},
		{
			name: "forbidden",
			mockBehaviour: func(s *mockService.MockTransactions, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, service.ErrAccountForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"account forbidden to access"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockTransactions, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, errors.New("general error"))
			},
			expectedCodeStatus:   500,
			expectedResponseBody: `{"message":"general error"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockTransactions, a *mockService.MockAccounts) {
				a.EXPECT().Get(context.Background(), accountID, userID).Return(domain.Account{}, nil)
				s.EXPECT().List(context.Background(), gomock.Any()).Return(transactions, errors.New("general error"))
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

			tService := mockService.NewMockTransactions(c)
			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(tService, aService)

			services := &service.Services{Transactions: tService, Accounts: aService}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/accounts/:id/transactions", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.listTransactionsOfAccount)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/accounts/%d/transactions", accountID),
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_createTransaction(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactions)

	date := time.Date(2021, time.September, 11, 21, 23, 21, 0, time.UTC)
	dateString := date.Format("2006-01-02T15:04:05.999999999Z07:00")

	toCreate := domain.TransactionToCreate{
		Amount:    100,
		Type:      domain.Income,
		CreatedAt: date,
	}

	expenseToCreate := domain.TransactionToCreate{
		Amount:    100,
		Type:      domain.Expense,
		CreatedAt: date,
	}

	transferToCreate := domain.TransactionToCreate{
		Amount:    100,
		Type:      domain.Transfer,
		CreatedAt: date,
	}

	created := domain.Transaction{
		ID:        1,
		Amount:    100,
		Type:      domain.Income,
		CreatedAt: time.Now(),
	}

	categoryId, creditId, debitId := new(int64), new(int64), new(int64)
	*categoryId, *creditId, *debitId = 1, 1, 2

	setResponseBody := func(t domain.Transaction) string {
		body, _ := json.Marshal(t)

		return string(body)
	}

	tests := []struct {
		name                 string
		requestBody          string
		categoryId           int64
		creditId             int64
		debitId              int64
		requestToCreate      domain.TransactionToCreate
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name:            "ok - income",
			requestBody:     fmt.Sprintf(`{"amount":100,"type":"income","createdAt":"%s"}`, dateString),
			categoryId:      *categoryId,
			debitId:         *debitId,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Create(context.Background(), toCreate, userID, categoryId, nil, debitId).Return(created, nil)
			},
			expectedCodeStatus:   201,
			expectedResponseBody: setResponseBody(created),
		},
		{
			name:            "ok - expanse",
			requestBody:     fmt.Sprintf(`{"amount":100,"type":"expense","createdAt":"%s"}`, dateString),
			categoryId:      *categoryId,
			creditId:        *creditId,
			requestToCreate: expenseToCreate,
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Create(context.Background(), expenseToCreate, userID, categoryId, creditId, nil).Return(created, nil)
			},
			expectedCodeStatus:   201,
			expectedResponseBody: setResponseBody(created),
		},
		{
			name:            "ok - transfer",
			requestBody:     fmt.Sprintf(`{"amount":100,"type":"transfer","createdAt":"%s"}`, dateString),
			categoryId:      *categoryId,
			debitId:         *debitId,
			creditId:        *creditId,
			requestToCreate: transferToCreate,
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Create(context.Background(), transferToCreate, userID, categoryId, creditId, debitId).Return(created, nil)
			},
			expectedCodeStatus:   201,
			expectedResponseBody: setResponseBody(created),
		},
		{
			name:                 "invalid request body",
			requestBody:          `{"amount":100,"type":"transfer"}`,
			mockBehaviour:        func(s *mockService.MockTransactions) {},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"invalid request body - Key: 'TransactionToCreate.CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag"}`,
		},
		{
			name:            "type mismatch",
			requestBody:     fmt.Sprintf(`{"amount":100,"type":"income","createdAt":"%s"}`, dateString),
			categoryId:      *categoryId,
			debitId:         *debitId,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Create(context.Background(), toCreate, userID, categoryId, nil, debitId).
					Return(created, service.ErrTransactionAndCategoryTypesMismatch)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"type of transaction and category does not match"}`,
		},
		{
			name:            "forbidden",
			requestBody:     fmt.Sprintf(`{"amount":100,"type":"income","createdAt":"%s"}`, dateString),
			categoryId:      *categoryId,
			debitId:         *debitId,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Create(context.Background(), toCreate, userID, categoryId, nil, debitId).
					Return(created, service.ErrDebitAccountForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"receiver account forbidden to access"}`,
		},
		{
			name:            "error",
			requestBody:     fmt.Sprintf(`{"amount":100,"type":"income","createdAt":"%s"}`, dateString),
			categoryId:      *categoryId,
			debitId:         *debitId,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Create(context.Background(), toCreate, userID, categoryId, nil, debitId).
					Return(created, errors.New("default error"))
			},
			expectedCodeStatus:   500,
			expectedResponseBody: `{"message":"default error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			tService := mockService.NewMockTransactions(c)
			tt.mockBehaviour(tService)

			services := &service.Services{Transactions: tService}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/transactions", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.createTransaction)

			// Create Request
			w := httptest.NewRecorder()
			queryString := "?"

			if tt.categoryId != 0 {
				queryString += fmt.Sprintf("categoryId=%d&", tt.categoryId)
			}

			if tt.creditId != 0 {
				queryString += fmt.Sprintf("creditId=%d&", tt.creditId)
			}

			if tt.debitId != 0 {
				queryString += fmt.Sprintf("debitId=%d&", tt.debitId)
			}

			req := httptest.NewRequest("POST", "/transactions"+queryString, bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
		})
	}
}

func TestHandler_deleteTransaction(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactions)

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Delete(context.Background(), transactionID, userID).Return(nil)
			},
			expectedCodeStatus:   204,
			expectedResponseBody: "",
		},
		{
			name: "access to account forbidden",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Delete(context.Background(), transactionID, userID).Return(service.ErrTransactionForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"transaction forbidden to access"}`,
		},
		{
			name: "account not found",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Delete(context.Background(), transactionID, userID).Return(repo.ErrTransactionNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"transaction doesn't exists"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockTransactions) {
				s.EXPECT().Delete(context.Background(), transactionID, userID).Return(errors.New("general error"))
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

			tService := mockService.NewMockTransactions(c)
			tt.mockBehaviour(tService)

			services := &service.Services{Transactions: tService}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/transactions/:id", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.deleteTransaction)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/transactions/%d", transactionID), bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

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
				s: services,
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

func TestHandler_listTransactionTypes(t *testing.T) {
	type mockBehaviour func(s *mockService.MockTransactionTypes)

	types := []domain.TransactionType{
		domain.Income,
	}

	setResponseBody := func(types []domain.TransactionType) string {
		body, _ := json.Marshal(types)

		return string(body)
	}

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockTransactionTypes) {
				s.EXPECT().List(context.Background()).Return(types, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(types),
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockTransactionTypes) {
				s.EXPECT().List(context.Background()).Return(types, errors.New("general error"))
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

			ttService := mockService.NewMockTransactionTypes(c)
			tt.mockBehaviour(ttService)

			services := &service.Services{TransactionTypes: ttService}
			handler := &Handler{
				s: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/transaction-types", handler.listTransactionTypes)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/transaction-types", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
