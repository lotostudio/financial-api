package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/service"
	mockService "github.com/lotostudio/financial-api/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_listUsers(t *testing.T) {
	type mockBehaviour func(s *mockService.MockUsers)

	users := []domain.User{
		{
			ID:        int64(1),
			FirstName: "Sirius",
			LastName:  "Sam",
			Email:     "qweqweqwe@gmail.com",
			Password:  "qweqweqwe",
		},
	}

	setResponseBody := func(users []domain.User) string {
		body, _ := json.Marshal(users)

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
			mockBehaviour: func(s *mockService.MockUsers) {
				s.EXPECT().List(context.Background()).Return(users, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: setResponseBody(users),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mockService.NewMockUsers(c)
			tt.mockBehaviour(userService)

			services := &service.Services{Users: userService}
			handler := &Handler{
				services: services,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/users", handler.listUsers)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedCodeStatus, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
