package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/service"
	"net/http"
	"strconv"
)

// @Summary Get statement
// @Tags statement
// @Description List transactions
// @ID getStatement
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
// @Failure 500 {object} response "Server error"
// @Router /stats/statement [get]
func (h *Handler) getStatement(c *gin.Context) {
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

	stat, err := h.s.Stats.Statement(c.Request.Context(), filter)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, stat)
}
