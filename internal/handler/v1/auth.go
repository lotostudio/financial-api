package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"net/http"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	users := api.Group("/auth")
	{
		users.POST("/register", h.register)
		users.POST("/login", h.login)
	}
}

// @Summary Register
// @Tags auth
// @Description User registration
// @ID register
// @Accept json
// @Produce json
// @Param input body domain.UserToCreate true "Register info"
// @Success 201 {string} domain.User "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 500 {object} response "Server error"
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var toCreate domain.UserToCreate

	if err := c.ShouldBindJSON(&toCreate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	user, err := h.services.Register(c.Request.Context(), toCreate)

	if err != nil {
		if err == repo.ErrUserAlreadyExists {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Login
// @Tags auth
// @Description User login
// @ID login
// @Accept json
// @Produce json
// @Param input body domain.UserToLogin true "Login credentials"
// @Success 200 {object} domain.Tokens "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 500 {object} response "Server error"
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var toLogin domain.UserToLogin

	if err := c.ShouldBindJSON(&toLogin); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	token, err := h.services.Login(c.Request.Context(), toLogin)

	if err != nil {
		if err == repo.ErrUserNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}
