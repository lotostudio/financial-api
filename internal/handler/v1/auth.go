package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/refresh", h.refresh)
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
// @Header 200 {int} Access-Token-TTL "Time to live of access token in seconds"
// @Header 200 {int} Refresh-Token-TTL "Time to live of refresh token in seconds"
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var toLogin domain.UserToLogin

	if err := c.ShouldBindJSON(&toLogin); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid request body - "+err.Error())
		return
	}

	tokens, err := h.services.Login(c.Request.Context(), toLogin)

	if err != nil {
		if err == repo.ErrUserNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Access-Token-TTL", strconv.Itoa(int(tokens.AccessTokenExpiredAt)))
	c.Header("Refresh-Token-TTL", strconv.Itoa(int(tokens.RefreshTokenExpiredAt)))

	c.JSON(http.StatusOK, tokens)
}

// @Summary Refresh
// @Tags auth
// @Description Refresh
// @ID refresh
// @Accept json
// @Produce json
// @Success 200 {object} domain.Tokens "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Header 200 {int} Access-Token-TTL "Time to live of access token in seconds"
// @Header 200 {int} Refresh-Token-TTL "Time to live of refresh token in seconds"
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := h.parseAuthHeaderToken(c)

	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	tokens, err := h.services.Refresh(c.Request.Context(), refreshToken)

	if err != nil {
		if err == repo.ErrSessionNotFound || err == service.ErrRefreshTokenExpired {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Access-Token-TTL", strconv.Itoa(int(tokens.AccessTokenExpiredAt)))
	c.Header("Refresh-Token-TTL", strconv.Itoa(int(tokens.RefreshTokenExpiredAt)))

	c.JSON(http.StatusOK, tokens)
}
