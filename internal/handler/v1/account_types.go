package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initAccountTypesRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/account-types")
	{
		accounts.GET("", h.listAccountTypes)
	}
}

// @Summary List account types
// @Tags accounts
// @Description List all account types
// @ID listAccountTypes
// @Accept json
// @Produce json
// @Success 200 {array} domain.AccountType "Operation finished successfully"
// @Failure 500 {object} response "Server error"
// @Router /accounts-types [get]
func (h *Handler) listAccountTypes(c *gin.Context) {
	accounts, err := h.services.AccountTypes.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}
