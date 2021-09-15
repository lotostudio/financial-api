package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/service"
	"net/http"
	"strconv"
	"time"
)

const (
	// Layout for dates in query params
	layout = "2006-01-02"
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
// @Param category query string false "Type of category"
// @Param dateFrom query string false "Start date (yyyy-MM-dd). Combined with dateTo"
// @Param dateTo query string false "End date (yyyy-MM-dd). Combined with dateFrom"
// @Success 200 {array} domain.Transaction "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /transactions [get]
func (h *Handler) listTransactions(c *gin.Context) {
	filter, err := h.parseTransactionsFilter(c)

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

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

	filter.OwnerId = &userId

	transactions, err := h.services.Transactions.List(c.Request.Context(), filter)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary List transactions of account
// @Tags transactions
// @Description List transactions of account
// @ID listTransactionsOfAccount
// @Accept json
// @Produce json
// @Param id path int true "Id of account"
// @Param category query string false "Type of category"
// @Param dateFrom query string false "Start date (yyyy-MM-dd)"
// @Param dateTo query string false "End date (yyyy-MM-dd)"
// @Success 200 {array} domain.Transaction "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /accounts/{id}/transactions [get]
func (h *Handler) listTransactionsOfAccount(c *gin.Context) {
	filter, err := h.parseTransactionsFilter(c)

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

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

	filter.OwnerId = &userId

	accountIdString := c.Param("id")

	if accountIdString == "" {
		newResponse(c, http.StatusBadRequest, "path param 'id' missing")
		return
	}

	accountId, err := strconv.ParseInt(accountIdString, 10, 32)

	if err != nil {
		newResponse(c, http.StatusBadRequest, "path param 'id' must be integer - "+err.Error())
		return
	}

	filter.AccountId = &accountId

	_, err = h.services.Accounts.Get(c.Request.Context(), accountId, userId)

	if err != nil {
		if err == service.ErrAccountForbidden {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		if err == repo.ErrAccountNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transactions, err := h.services.Transactions.List(c.Request.Context(), filter)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Parse query params to domain.TransactionsFilter
func (h *Handler) parseTransactionsFilter(c *gin.Context) (domain.TransactionsFilter, error) {
	filter := domain.TransactionsFilter{}

	category := c.Query("category")

	if category != "" {
		filter.Category = &category
	}

	dateFromString := c.Query("dateFrom")

	if dateFromString != "" {
		dateFrom, err := time.Parse(layout, dateFromString)

		if err != nil {
			return filter, err
		}

		filter.CreatedFrom = &dateFrom
	}

	dateToString := c.Query("dateTo")

	if dateToString != "" {
		dateTo, err := time.Parse(layout, dateToString)

		if err != nil {
			return filter, err
		}

		filter.CreatedTo = &dateTo
	}

	// Filters by date period both must be null or not null at the same time
	if (filter.CreatedFrom != nil && filter.CreatedTo == nil) || (filter.CreatedFrom == nil && filter.CreatedTo != nil) {
		return filter, errDateFiltersInvalid
	}

	return filter, nil
}

// @Summary Create transaction
// @Tags transactions
// @Description Create new transaction by type:
// @Description * income - pass debit account and transaction category
// @Description * expense - pass credit account and transaction category
// @Description * transfer - pass credit and debit accounts
// @ID createTransaction
// @Accept json
// @Produce json
// @Param categoryId query int false "Id of category"
// @Param creditId query int false "Id of credit account"
// @Param debitId query int false "Id of debit account"
// @Param input body domain.TransactionToCreate true "Transaction info"
// @Success 201 {array} domain.Transaction "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
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
		if err == service.ErrTransactionAndCategoryTypesMismatch || err == repo.ErrTransactionCategoryNotFound ||
			err == service.ErrNoAccountSelected || err == service.ErrAccountsHaveDifferenceCurrencies ||
			err == repo.ErrAccountNotEnoughBalance {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err == service.ErrDebitAccountForbidden || err == service.ErrCreditAccountForbidden {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
