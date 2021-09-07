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
)

const (
	userID     = int64(1)
	accountID  = int64(2)
	currencyID = 3
)

func TestHandler_listAccounts(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAccounts)

	accounts := []domain.Account{
		{
			ID:       1,
			Title:    "acc1",
			Balance:  12.1,
			Currency: "KZT",
			Type:     domain.Card,
		},
	}

	setResponseBody := func(accounts []domain.Account) string {
		body, _ := json.Marshal(accounts)

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
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().List(context.Background(), userID).Return(accounts, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(accounts),
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().List(context.Background(), userID).Return(accounts, errors.New("general error"))
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

			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(aService)

			services := &service.Services{Accounts: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/accounts", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.listAccounts)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/accounts", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_createAccount(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAccounts)

	currencyIDString := strconv.FormatInt(currencyID, 10)

	toCreate := domain.AccountToCreate{
		Title:   "Acc1",
		Balance: 12,
		Type:    domain.Card,
	}

	account := domain.Account{
		ID:       1,
		Title:    "Acc1",
		Balance:  12,
		Currency: "KZT",
		Type:     domain.Card,
		OwnerId:  userID,
	}

	setResponseBody := func(account domain.Account) string {
		body, _ := json.Marshal(account)

		return string(body)
	}

	tests := []struct {
		name                 string
		currencyId           string
		requestBody          string
		requestToCreate      domain.AccountToCreate
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name:            "ok",
			currencyId:      currencyIDString,
			requestBody:     `{"title":"Acc1","balance":12,"type":"card"}`,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Create(context.Background(), toCreate, userID, currencyID).Return(account, nil)
			},
			expectedCodeStatus:   201,
			expectedResponseBody: setResponseBody(account),
		},
		{
			name:                 "missing currency id",
			currencyId:           "",
			mockBehaviour:        func(s *mockService.MockAccounts) {},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"query param 'currencyId' missing"}`,
		},
		{
			name:                 "currency id not integer",
			currencyId:           "qwe",
			mockBehaviour:        func(s *mockService.MockAccounts) {},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"query param 'currencyId' must be integer - strconv.ParseInt: parsing \"qwe\": invalid syntax"}`,
		},
		{
			name:                 "invalid request body",
			currencyId:           currencyIDString,
			requestBody:          `{"title":"Acc1","balance":12,"type":"qwe"}`,
			mockBehaviour:        func(s *mockService.MockAccounts) {},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"invalid request body - Key: 'AccountToCreate.Type' Error:Field validation for 'Type' failed on the 'oneof' tag"}`,
		},
		{
			name:            "currency not found",
			currencyId:      currencyIDString,
			requestBody:     `{"title":"Acc1","balance":12,"type":"card"}`,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Create(context.Background(), toCreate, userID, currencyID).Return(account, repo.ErrCurrencyNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"currency doesn't exists"}`,
		},
		{
			name:        "invalid loan data",
			currencyId:  currencyIDString,
			requestBody: `{"title":"Acc1","balance":12,"type":"loan"}`,
			requestToCreate: domain.AccountToCreate{
				Title:   "Acc1",
				Balance: 12,
				Type:    domain.Loan,
				Term:    nil,
				Rate:    nil,
			},
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Create(context.Background(), domain.AccountToCreate{
					Title:   "Acc1",
					Balance: 12,
					Type:    domain.Loan,
					Term:    nil,
					Rate:    nil,
				}, userID, currencyID).Return(account, service.ErrInvalidLoanData)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account with type 'loan' must have valid term and rate"}`,
		},
		{
			name:        "invalid deposit data",
			currencyId:  currencyIDString,
			requestBody: `{"title":"Acc1","balance":12,"type":"deposit"}`,
			requestToCreate: domain.AccountToCreate{
				Title:   "Acc1",
				Balance: 12,
				Type:    domain.Deposit,
				Term:    nil,
				Rate:    nil,
			},
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Create(context.Background(), domain.AccountToCreate{
					Title:   "Acc1",
					Balance: 12,
					Type:    domain.Deposit,
					Term:    nil,
					Rate:    nil,
				}, userID, currencyID).Return(account, service.ErrInvalidDepositData)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account with type 'deposit' must have valid term and rate"}`,
		},
		{
			name:            "error",
			currencyId:      currencyIDString,
			requestBody:     `{"title":"Acc1","balance":12,"type":"card"}`,
			requestToCreate: toCreate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Create(context.Background(), toCreate, userID, currencyID).Return(account, errors.New("general error"))
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

			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(aService)

			services := &service.Services{Accounts: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/accounts", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.createAccount)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/accounts?currencyId="+tt.currencyId, bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAccount(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAccounts)

	account := domain.Account{
		ID:       1,
		Title:    "Acc1",
		Balance:  12.1,
		Currency: "KZT",
		Type:     domain.Card,
	}

	setResponseBody := func(account domain.Account) string {
		body, _ := json.Marshal(account)

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
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Get(context.Background(), accountID, userID).Return(account, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(account),
		},
		{
			name: "access to account forbidden",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Get(context.Background(), accountID, userID).Return(account, service.ErrAccountForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"account forbidden to access"}`,
		},
		{
			name: "account not found",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Get(context.Background(), accountID, userID).Return(account, repo.ErrAccountNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account doesn't exists"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Get(context.Background(), accountID, userID).Return(account, errors.New("general error"))
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

			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(aService)

			services := &service.Services{Accounts: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/accounts/:id", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.getAccount)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/accounts/%d", accountID), bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateAccount(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAccounts)

	title, balance := "Acc1", float64(12)
	toUpdate := domain.AccountToUpdate{
		Title:   &title,
		Balance: &balance,
	}

	account := domain.Account{
		ID:       1,
		Title:    "Acc1",
		Balance:  12,
		Currency: "KZT",
		Type:     domain.Card,
		OwnerId:  userID,
	}

	setResponseBody := func(account domain.Account) string {
		body, _ := json.Marshal(account)

		return string(body)
	}

	tests := []struct {
		name                 string
		requestBody          string
		requestToUpdate      domain.AccountToUpdate
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name:            "ok",
			requestBody:     `{"title":"Acc1","balance":12}`,
			requestToUpdate: toUpdate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Update(context.Background(), toUpdate, accountID, userID).Return(account, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(account),
		},
		{
			name:                 "invalid request body",
			requestBody:          `{"title":"Acc1","balance":-1}`,
			mockBehaviour:        func(s *mockService.MockAccounts) {},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"invalid request body - Key: 'AccountToUpdate.Balance' Error:Field validation for 'Balance' failed on the 'gte' tag"}`,
		},
		{
			name:            "access to account forbidden",
			requestBody:     `{"title":"Acc1","balance":12}`,
			requestToUpdate: toUpdate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Update(context.Background(), toUpdate, accountID, userID).Return(account, service.ErrAccountForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"account forbidden to access"}`,
		},
		{
			name:            "account not found",
			requestBody:     `{"title":"Acc1","balance":12}`,
			requestToUpdate: toUpdate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Update(context.Background(), toUpdate, accountID, userID).Return(account, repo.ErrAccountNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account doesn't exists"}`,
		},
		{
			name:        "invalid loan data",
			requestBody: `{"title":"Acc1","balance":12}`,
			requestToUpdate: domain.AccountToUpdate{
				Title:   &title,
				Balance: &balance,
			},
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Update(context.Background(), domain.AccountToUpdate{
					Title:   &title,
					Balance: &balance,
				}, accountID, userID).Return(account, service.ErrInvalidLoanData)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account with type 'loan' must have valid term and rate"}`,
		},
		{
			name:        "invalid loan data",
			requestBody: `{"title":"Acc1","balance":12}`,
			requestToUpdate: domain.AccountToUpdate{
				Title:   &title,
				Balance: &balance,
			},
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Update(context.Background(), domain.AccountToUpdate{
					Title:   &title,
					Balance: &balance,
				}, accountID, userID).Return(account, service.ErrInvalidDepositData)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account with type 'deposit' must have valid term and rate"}`,
		},
		{
			name:            "error",
			requestBody:     `{"title":"Acc1","balance":12}`,
			requestToUpdate: toUpdate,
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Update(context.Background(), toUpdate, accountID, userID).Return(account, errors.New("general error"))
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

			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(aService)

			services := &service.Services{Accounts: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.PUT("/accounts/:id", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.updateAccount)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/accounts/%d", accountID), bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteAccount(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAccounts)

	tests := []struct {
		name                 string
		mockBehaviour        mockBehaviour
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Delete(context.Background(), accountID, userID).Return(nil)
			},
			expectedCodeStatus:   204,
			expectedResponseBody: "",
		},
		{
			name: "access to account forbidden",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Delete(context.Background(), accountID, userID).Return(service.ErrAccountForbidden)
			},
			expectedCodeStatus:   403,
			expectedResponseBody: `{"message":"account forbidden to access"}`,
		},
		{
			name: "account not found",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Delete(context.Background(), accountID, userID).Return(repo.ErrAccountNotFound)
			},
			expectedCodeStatus:   400,
			expectedResponseBody: `{"message":"account doesn't exists"}`,
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockAccounts) {
				s.EXPECT().Delete(context.Background(), accountID, userID).Return(errors.New("general error"))
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

			aService := mockService.NewMockAccounts(c)
			tt.mockBehaviour(aService)

			services := &service.Services{Accounts: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/accounts/:id", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.deleteAccount)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/accounts/%d", accountID), bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
