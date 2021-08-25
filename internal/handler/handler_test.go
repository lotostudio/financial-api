package handler

import (
	"github.com/lotostudio/financial-api/internal/config"
	"github.com/lotostudio/financial-api/internal/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(&service.Services{})

	require.IsType(t, &Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {
	h := NewHandler(&service.Services{})

	router := h.Init(&config.Config{})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")

	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
