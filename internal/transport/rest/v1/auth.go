package v1

import (
	"errors"
	"net/http"
	"seal/pkg/app_error"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerAuthHandler(api *gin.RouterGroup) {
	group := api.Group("/auth")
	{
		group.POST("sign-in", h.signIn)
		group.POST("refresh", h.refresh)
	}
}

type Auth struct {
	Login    string `json:"login" validate:"required,max=50,min=5"`
	Password string `json:"password" validate:"required,min=6"`
}

type refreshTokenBody struct {
	RefreshToken string `json:"refresh_token" validate:"required,max=500,min=10"`
}

// Login godoc
// @Summary      Login
// @Description  login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param		 data	body	Auth	true	"data"
// @Success      200	{object}	transport.JwtWithRefresh
// @Failure      400 	{object}	app_error.AppError
// @Failure      403 	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var userFromQuery Auth
	if err := c.ShouldBind(&userFromQuery); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(userFromQuery); errs != nil {
		c.Error(app_error.ValidationError(errs))
		return
	}

	user, err := h.Usecase.User.GetByCredentials(userFromQuery.Login, userFromQuery.Password)

	if errors.Is(err, app_error.ErrNotFound) {
		c.Error(app_error.LoginPasswordError())
		return
	} else if err != nil {
		c.Error(err)
		return
	}

	if jwt, err := h.JwtWorker.GenerateTokens(user.Id, user.Role, user.Title); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, jwt)
	}
}

// RefreshToken godoc
// @Summary      Refresh tokens
// @Description  refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param		 data	body	refreshTokenBody	true	"data"
// @Success      200 	{object}	transport.JwtWithRefresh
// @Failure      401	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /auth/refresh [post]
// @Security 	 BearerAuth
func (h *Handler) refresh(c *gin.Context) {
	//token, err := middleware.GetBearerToken(c.GetHeader("Authorization"))

	var token refreshTokenBody
	if err := c.ShouldBind(&token); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(token); errs != nil {
		c.Error(app_error.ValidationError(errs))
		return
	}

	userId, err := h.JwtWorker.GetUserIdFromRefreshToken(token.RefreshToken)

	if err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.ErrUnauthorized)
		return
	}

	user, err := h.Usecase.User.GetById(userId)

	if errors.Is(err, app_error.ErrNotFound) {
		c.Error(app_error.ErrUnauthorized)
		return
	} else if err != nil {
		c.Error(err)
		return
	}

	jwt, err := h.JwtWorker.GenerateTokens(user.Id, user.Role, user.Title)

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, jwt)
	}

}
