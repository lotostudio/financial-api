package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/service"
	"github.com/lotostudio/financial-api/pkg/auth"
)

type Handler struct {
	s   *service.Services
	tkn auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		s:   services,
		tkn: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initAuthRoutes(v1)
		h.initAccountsRoutes(v1)
		h.initTransactionsRoutes(v1)
	}
}
