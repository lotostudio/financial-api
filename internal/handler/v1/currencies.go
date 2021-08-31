package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initCurrenciesRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/currencies")
	{
		accounts.GET("", h.listCurrencies)
	}
}

// @Summary List currencies
// @Tags currencies
// @Description List all currencies
// @ID listCurrencies
// @Accept json
// @Produce json
// @Success 200 {array} domain.Currency "Operation finished successfully"
// @Failure 500 {object} response "Server error"
// @Router /currencies [get]
func (h *Handler) listCurrencies(c *gin.Context) {
	accounts, err := h.services.Currencies.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}
