package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"net/http"
	"strconv"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.listUsers)

		me := users.Group("/me", h.userIdentity)
		{
			me.PATCH("", h.partialUpdateMe)
		}
	}
}

// @Summary List users
// @Tags users
// @Description List all users
// @ID listUsers
// @Accept json
// @Produce json
// @Success 200 {array} domain.User "Operation finished successfully"
// @Failure 500 {object} response "Server error"
// @Router /users [get]
func (h *Handler) listUsers(c *gin.Context) {
	users, err := h.services.Users.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Update me
// @Tags users
// @Description Update self user
// @ID partialUpdateMe
// @Accept json
// @Produce json
// @Param input body domain.UserToUpdate true "Data to update"
// @Success 200 {object} domain.User "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 500 {object} response "Server error"
// @Router /users/me [patch]
func (h *Handler) partialUpdateMe(c *gin.Context) {
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

	var toUpdate domain.UserToUpdate

	if err = c.ShouldBindJSON(&toUpdate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	user, err := h.services.Users.UpdatePassword(c.Request.Context(), userId, toUpdate)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
