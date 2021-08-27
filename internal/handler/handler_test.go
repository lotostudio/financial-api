package handler

import (
	"github.com/lotostudio/financial-api/internal/config"
	"github.com/lotostudio/financial-api/internal/service"
	"github.com/lotostudio/financial-api/pkg/auth"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	tokenManager, _ := auth.NewJWTManager("key", 5*time.Second)

	h := NewHandler(&service.Services{}, tokenManager)

	require.IsType(t, &Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {
	tokenManager, _ := auth.NewJWTManager("key", 5*time.Second)

	h := NewHandler(&service.Services{}, tokenManager)

	router := h.Init(&config.Config{})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")

	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
