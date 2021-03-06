package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) initAccountsRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/accounts", h.userIdentity)
	{
		accounts.GET("", h.listAccounts)
		accounts.GET("/grouped", h.listGropedAccounts)
		accounts.POST("", h.createAccount)

		account := accounts.Group("/:id")
		{
			account.GET("", h.getAccount)
			account.PUT("", h.updateAccount)
			account.DELETE("", h.deleteAccount)

			statement := account.Group("/statement")
			{
				statement.GET("", h.getStatement)
			}

			transactions := account.Group("/transactions")
			{
				transactions.GET("", h.listTransactionsOfAccount)
			}
		}
	}

	types := api.Group("/account-types")
	{
		types.GET("", h.listAccountTypes)
	}

	currencies := api.Group("/currencies")
	{
		currencies.GET("", h.listCurrencies)
	}
}

// @Summary List accounts
// @Tags accounts
// @Description List accounts of user
// @ID listAccounts
// @Security UsersAuth
// @Accept json
// @Produce json
// @Success 200 {array} domain.Account "Operation finished successfully"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /accounts [get]
func (h *Handler) listAccounts(c *gin.Context) {
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

	accounts, err := h.s.Accounts.List(c.Request.Context(), userId)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// @Summary List grouped accounts
// @Tags accounts
// @Description List grouped accounts of user by types
// @ID listGroupedAccounts
// @Security UsersAuth
// @Accept json
// @Produce json
// @Success 200 {object} domain.GroupedAccounts "Operation finished successfully"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /accounts/grouped [get]
func (h *Handler) listGropedAccounts(c *gin.Context) {
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

	accounts, err := h.s.Accounts.ListGrouped(c.Request.Context(), userId)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// @Summary Create account
// @Tags accounts
// @Description Create new account
// @ID createAccount
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param currencyId query int true "Id of currency"
// @Param input body domain.AccountToCreate true "Account info"
// @Success 201 {object} domain.Account "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /accounts [post]
func (h *Handler) createAccount(c *gin.Context) {
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

	currencyIdString := c.Query("currencyId")

	if currencyIdString == "" {
		newResponse(c, http.StatusBadRequest, "query param 'currencyId' missing")
		return
	}

	currencyId, err := strconv.ParseInt(currencyIdString, 32, 10)

	if err != nil {
		newResponse(c, http.StatusBadRequest, "query param 'currencyId' must be integer - "+err.Error())
		return
	}

	var toCreate domain.AccountToCreate

	if err := c.ShouldBindJSON(&toCreate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	account, err := h.s.Accounts.Create(c.Request.Context(), toCreate, userId, int(currencyId))

	if errors.Is(err, repo.ErrCurrencyNotFound) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrInvalidLoanData) || errors.Is(err, service.ErrInvalidDepositData) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrAccountCountLimited) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

// @Summary Get account
// @Tags accounts
// @Description Get account of user
// @ID getAccount
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param id path int64 true "Id of account"
// @Success 200 {object} domain.Account "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /accounts/{id} [get]
func (h *Handler) getAccount(c *gin.Context) {
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

	id, err := strconv.ParseInt(idString, 10, 32)

	if err != nil {
		newResponse(c, http.StatusBadRequest, "path param 'id' must be integer - "+err.Error())
		return
	}

	accounts, err := h.s.Accounts.Get(c.Request.Context(), id, userId)

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

	c.JSON(http.StatusOK, accounts)
}

// @Summary Update account
// @Tags accounts
// @Description Update account of user
// @ID updateAccount
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param id path int64 true "Id of account"
// @Param input body domain.AccountToUpdate true "Account info"
// @Success 200 {object} domain.Account "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /accounts/{id} [put]
func (h *Handler) updateAccount(c *gin.Context) {
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

	id, err := strconv.ParseInt(idString, 10, 32)

	if err != nil {
		newResponse(c, http.StatusBadRequest, "path param 'id' must be integer - "+err.Error())
		return
	}

	var toUpdate domain.AccountToUpdate

	if err := c.ShouldBindJSON(&toUpdate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	accounts, err := h.s.Accounts.Update(c.Request.Context(), toUpdate, id, userId)

	if errors.Is(err, service.ErrAccountForbidden) {
		newResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if errors.Is(err, repo.ErrAccountNotFound) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, service.ErrInvalidLoanData) || errors.Is(err, service.ErrInvalidDepositData) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// @Summary Delete account
// @Tags accounts
// @Description Delete account of user
// @ID deleteAccount
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param id path int64 true "Id of account"
// @Success 204 {null} nil "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /accounts/{id} [delete]
func (h *Handler) deleteAccount(c *gin.Context) {
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

	id, err := strconv.ParseInt(idString, 10, 32)

	if err != nil {
		newResponse(c, http.StatusBadRequest, "path param 'id' must be integer - "+err.Error())
		return
	}

	err = h.s.Accounts.Delete(c.Request.Context(), id, userId)

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

	c.Status(http.StatusNoContent)
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
	accounts, err := h.s.AccountTypes.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
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
	accounts, err := h.s.Currencies.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}
