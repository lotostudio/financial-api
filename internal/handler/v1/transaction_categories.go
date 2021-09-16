package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"net/http"
)

func (h *Handler) initTransactionCategoriesRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/transaction-categories")
	{
		accounts.GET("", h.listTransactionCategories)
	}
}

// @Summary List transaction categories
// @Tags transactions
// @Description List transaction categories
// @ID listTransactionCategories
// @Accept json
// @Produce json
// @Param type query string false "Type of category"
// @Success 200 {array} domain.TransactionCategory "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 500 {object} response "Server error"
// @Router /transaction-categories [get]
func (h *Handler) listTransactionCategories(c *gin.Context) {
	_type := domain.TransactionType(c.Query("type"))

	if _type == "" {
		accounts, err := h.services.TransactionCategories.List(c.Request.Context())

		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, accounts)
		return
	}

	accounts, err := h.services.TransactionCategories.ListByType(c.Request.Context(), _type)

	if err != nil {
		if err == domain.ErrInvalidTransactionType {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}
