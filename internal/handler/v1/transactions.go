package v1

import (
	"errors"
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
	transactions := api.Group("/transactions", h.userIdentity)
	{
		transactions.GET("", h.listTransactions)
		transactions.POST("", h.createTransaction)
		transactions.DELETE("/:id", h.deleteTransaction)

		stats := transactions.Group("/stats")
		{
			stats.GET("", h.transactionStats)
		}
	}

	categories := api.Group("/transaction-categories")
	{
		categories.GET("", h.listTransactionCategories)
	}

	types := api.Group("/transaction-types")
	{
		types.GET("", h.listTransactionTypes)
	}
}

// @Summary List transactions
// @Tags transactions
// @Description List transactions
// @ID listTransactions
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param category query string false "Type of transaction"
// @Param type query string false "Type of transaction"
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

	transactions, err := h.s.Transactions.List(c.Request.Context(), filter)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary Transaction stats
// @Tags transactions
// @Description Transaction stats
// @ID transactionStats
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param category query string false "Type of transaction"
// @Param type query string false "Type of transaction"
// @Param dateFrom query string false "Start date (yyyy-MM-dd). Combined with dateTo"
// @Param dateTo query string false "End date (yyyy-MM-dd). Combined with dateFrom"
// @Success 200 {array} domain.TransactionStat "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /transactions/stats [get]
func (h *Handler) transactionStats(c *gin.Context) {
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

	stats, err := h.s.Transactions.Stats(c.Request.Context(), filter)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary List transactions of account
// @Tags transactions
// @Description List transactions of account
// @ID listTransactionsOfAccount
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param id path int true "Id of account"
// @Param category query string false "Category of transaction"
// @Param type query string false "Type of transaction"
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

	_, err = h.s.Accounts.Get(c.Request.Context(), accountId, userId)

	if errors.Is(err, service.ErrAccountForbidden) {
		newResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if errors.Is(err, repo.ErrAccountNotFound) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transactions, err := h.s.Transactions.List(c.Request.Context(), filter)

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

	_type := domain.TransactionType(c.Query("type"))

	if err := _type.Validate(); _type != "" && err != nil {
		return filter, err
	}

	filter.Type = &_type

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
// @Security UsersAuth
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

	var categoryID = new(int64)
	if categoryIDString := c.Query("categoryId"); categoryIDString != "" {
		*categoryID, err = strconv.ParseInt(categoryIDString, 10, 64)

		if err != nil {
			newResponse(c, http.StatusBadRequest, "query param 'creditId' must be integer - "+err.Error())
			return
		}
	} else {
		categoryID = nil
	}

	var creditID = new(int64)
	if creditIDString := c.Query("creditId"); creditIDString != "" {
		*creditID, err = strconv.ParseInt(creditIDString, 10, 64)

		if err != nil {
			newResponse(c, http.StatusBadRequest, "query param 'creditId' must be integer - "+err.Error())
			return
		}
	} else {
		creditID = nil
	}

	var debitID = new(int64)
	if debitIDString := c.Query("debitId"); debitIDString != "" {
		*debitID, err = strconv.ParseInt(debitIDString, 10, 64)

		if err != nil {
			newResponse(c, http.StatusBadRequest, "query param 'debitId' must be integer - "+err.Error())
			return
		}
	} else {
		debitID = nil
	}

	var toCreate domain.TransactionToCreate

	if err = c.ShouldBindJSON(&toCreate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	transaction, err := h.s.Transactions.Create(c.Request.Context(), toCreate, userId, categoryID, creditID, debitID)

	if errors.Is(err, service.ErrTransactionAndCategoryTypesMismatch) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, repo.ErrTransactionCategoryNotFound) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrNoAccountSelected) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrAccountsHaveDifferenceCurrencies) || errors.Is(err, repo.ErrAccountNotEnoughBalance) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrDebitAccountForbidden) || errors.Is(err, service.ErrCreditAccountForbidden) {
		newResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// @Summary Delete transaction
// @Tags transactions
// @Description Delete transaction
// @ID deleteTransaction
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param id path int64 true "Id of transaction"
// @Success 204 {null} nil "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /transactions/{id} [delete]
func (h *Handler) deleteTransaction(c *gin.Context) {
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

	idString := c.Param("id")

	if idString == "" {
		newResponse(c, http.StatusBadRequest, "path param 'id' missing")
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)

	if err != nil {
		newResponse(c, http.StatusBadRequest, "path param 'id' must be integer - "+err.Error())
		return
	}

	if err = h.s.Transactions.Delete(c.Request.Context(), id, userId); err != nil {
		if err == repo.ErrTransactionNotFound || err == repo.ErrAccountNotEnoughBalance {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err == service.ErrTransactionForbidden {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
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
		accounts, err := h.s.TransactionCategories.List(c.Request.Context())

		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, accounts)
		return
	}

	accounts, err := h.s.TransactionCategories.ListByType(c.Request.Context(), _type)

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
	types, err := h.s.TransactionTypes.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types)
}
