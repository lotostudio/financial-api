package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"net/http"
	"strconv"
)

func (h *Handler) initTransactionsRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/transactions", h.userIdentity)
	{
		accounts.GET("", h.listTransactions)
		accounts.POST("", h.createTransaction)
	}
}

// @Summary List transactions
// @Tags transactions
// @Description List transactions
// @ID listTransactions
// @Accept json
// @Produce json
// @Param type query string false "Type of category"
// @Success 200 {array} domain.Transaction "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /transactions [get]
func (h *Handler) listTransactions(c *gin.Context) {
	userIdString, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := strconv.ParseInt(userIdString.(string), 10, 64)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	transactions, err := h.services.Transactions.List(c.Request.Context(), userId)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary Create transaction
// @Tags transactions
// @Description Create transaction
// @ID createTransaction
// @Accept json
// @Produce json
// @Param type query string false "Type of category"
// @Param input body domain.TransactionToCreate true "Transaction info"
// @Success 200 {array} domain.Transaction "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /transactions [post]
func (h *Handler) createTransaction(c *gin.Context) {
	userIdString, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := strconv.ParseInt(userIdString.(string), 10, 64)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	var categoryId = new(int64)
	categoryIdString := c.Query("categoryId")

	if categoryIdString == "" {
		categoryId = nil
	} else {
		*categoryId, err = strconv.ParseInt(categoryIdString, 10, 64)

		if err != nil {
			newResponse(c, http.StatusBadRequest, "query param 'creditId' must be integer - "+err.Error())
			return
		}
	}

	var creditId = new(int64)
	creditIdString := c.Query("creditId")

	if creditIdString == "" {
		creditId = nil
	} else {
		*creditId, err = strconv.ParseInt(creditIdString, 10, 64)

		if err != nil {
			newResponse(c, http.StatusBadRequest, "query param 'creditId' must be integer - "+err.Error())
			return
		}
	}

	var debitId = new(int64)
	debitIdString := c.Query("debitId")

	if debitIdString == "" {
		debitId = nil
	} else {
		*debitId, err = strconv.ParseInt(debitIdString, 10, 64)

		if err != nil {
			newResponse(c, http.StatusBadRequest, "query param 'debitId' must be integer - "+err.Error())
			return
		}
	}

	var toCreate domain.TransactionToCreate

	if err := c.ShouldBindJSON(&toCreate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	transaction, err := h.services.Transactions.Create(c.Request.Context(), toCreate, userId, categoryId, creditId, debitId)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
