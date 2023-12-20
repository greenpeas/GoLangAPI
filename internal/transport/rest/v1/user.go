package v1

import (
	"net/http"
	"seal/internal/domain/user"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type userList struct {
	RecordsTotal    int         `json:"records_total"`
	RecordsFiltered int         `json:"records_filtered"`
	Data            []user.User `json:"data"`
}

func (h *Handler) registerUserHandler(api *gin.RouterGroup) {
	group := api.Group("/user")
	{
		group.GET(":id", h.user)
		group.PUT(":id", h.userUpdate)
		group.DELETE(":id", h.userDelete)
		group.GET("", h.userList)
		group.POST("", h.userCreate)
	}
}

// ItemUser godoc
// @Summary      User info
// @Description  user info
// @Tags         user
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	user.User
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /user/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) user(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.User.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListUsers godoc
// @Summary      List users
// @Description  get users
// @Tags         user
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	userList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /user [get]
// @Security 	 BearerAuth
func (h *Handler) userList(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.User.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// CreateUSer godoc
// @Summary      Create user
// @Description  add user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		 data	body	user.Db	true	"data"
// @Success      200	{object}	user.User
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /user [post]
// @Security 	 BearerAuth
func (h *Handler) userCreate(c *gin.Context) {
	var fromRequest user.CreateRequest
	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	if createdUser, err := h.Usecase.User.Create(fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdUser)
	}
}

// CreateUSer godoc
// @Summary      Update user
// @Description  update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	user.Db	true	"data"
// @Success      200	{object}	user.User
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /user/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) userUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest user.UpdateRequest

	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	if updatedUser, err := h.Usecase.User.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, updatedUser)
	}
}

// ItemRoute godoc
// @Summary      User delete
// @Description  User delete
// @Tags         user
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.DeleteResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /user/{id} [delete]
// @Security 	 BearerAuth
func (h *Handler) userDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.User.DeleteById(id); err != nil {
		c.Error(err)
	} else if data {
		c.JSON(http.StatusOK, transport.DeleteResponse{Success: true})
	} else {
		c.Error(app_error.ErrNotFound)
	}
}
