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
	"strconv"
	"testing"
)

func TestHandler_listCurrencies(t *testing.T) {
	type mockBehaviour func(s *mockService.MockCurrencies)

	currencies := []domain.Currency{
		{
			ID:   1,
			Code: "KZT",
		},
	}

	setResponseBody := func(currencies []domain.Currency) string {
		body, _ := json.Marshal(currencies)

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
			mockBehaviour: func(s *mockService.MockCurrencies) {
				s.EXPECT().List(context.Background()).Return(currencies, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(currencies),
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockCurrencies) {
				s.EXPECT().List(context.Background()).Return(currencies, errors.New("general error"))
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

			cService := mockService.NewMockCurrencies(c)
			tt.mockBehaviour(cService)

			services := &service.Services{Currencies: cService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/currencies", func(c *gin.Context) {
				c.Set(userCtx, strconv.FormatInt(userID, 10))
			}, handler.listCurrencies)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/currencies", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
