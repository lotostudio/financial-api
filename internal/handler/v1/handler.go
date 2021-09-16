package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/service"
	"github.com/lotostudio/financial-api/pkg/auth"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initAuthRoutes(v1)
		h.initCurrenciesRoutes(v1)
		h.initAccountsRoutes(v1)
		h.initAccountTypesRoutes(v1)
		h.initTransactionsRoutes(v1)
		h.initTransactionCategoriesRoutes(v1)
		h.initTransactionTypesRoutes(v1)
	}
}
