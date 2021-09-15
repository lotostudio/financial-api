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

func TestHandler_listAccountTypes(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAccountTypes)

	types := []domain.AccountType{
		domain.Card,
	}

	setResponseBody := func(currencies []domain.AccountType) string {
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
			mockBehaviour: func(s *mockService.MockAccountTypes) {
				s.EXPECT().List(context.Background()).Return(types, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(types),
		},
		{
			name: "error",
			mockBehaviour: func(s *mockService.MockAccountTypes) {
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

			aService := mockService.NewMockAccountTypes(c)
			tt.mockBehaviour(aService)

			services := &service.Services{AccountTypes: aService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/account-types", handler.listAccountTypes)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/account-types", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}