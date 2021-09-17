package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initTransactionTypesRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/transaction-types")
	{
		accounts.GET("", h.listTransactionTypes)
	}
}

// @Summary List transaction types
// @Tags transactions
// @Description List all transaction types
// @ID listTransactionTypes
// @Accept json
// @Produce json
// @Success 200 {array} domain.TransactionType "Operation finished successfully"
// @Failure 500 {object} response "Server error"
// @Router /transaction-types [get]
func (h *Handler) listTransactionTypes(c *gin.Context) {
	types, err := h.services.TransactionTypes.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types)
}
