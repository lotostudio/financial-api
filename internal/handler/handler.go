package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/config"
	v1 "github.com/lotostudio/financial-api/internal/handler/v1"
	"github.com/lotostudio/financial-api/internal/service"
	"github.com/lotostudio/financial-api/pkg/auth"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	// Init swagger routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Init router
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
