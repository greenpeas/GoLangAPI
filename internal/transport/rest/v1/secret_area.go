package v1

import (
	"errors"
	"net/http"
	"seal/internal/domain/secret_area"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type secretAreaList struct {
	RecordsTotal    int                      `json:"records_total"`
	RecordsFiltered int                      `json:"records_filtered"`
	Data            []secret_area.SecretArea `json:"data"`
}

func (h *Handler) registerSecretAreaHandler(api *gin.RouterGroup) {
	group := api.Group("/secret-area")
	{
		group.PUT(":id", h.secretAreaUpdate)
		group.GET(":id", h.secretArea)
		group.DELETE(":id", h.secretAreaDelete)
		group.POST("", h.createSecretArea)
		group.GET("", h.secretAreaList)
	}
}

// ItemSecretArea godoc
// @Summary      SecretArea info
// @Description  secretArea info
// @Tags         secret-area
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	secret_area.SecretArea
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /secret-area/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) secretArea(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.SecretArea.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListSecretArea godoc
// @Summary      List secret area
// @Description  get secret areas
// @Tags         secret-area
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	secretAreaList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /secret-area [get]
// @Security 	 BearerAuth
func (h *Handler) secretAreaList(c *gin.Context) {
	var queryParams secret_area.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.SecretArea.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// CreateSecretArea godoc
// @Summary      Create secret area
// @Description  add secret area
// @Tags         secret-area
// @Accept       json
// @Produce      json
// @Param		 data	body	secret_area.CreateRequest	true	"data"
// @Success      200	{object}	secret_area.SecretArea
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /secret-area [post]
// @Security 	 BearerAuth
func (h *Handler) createSecretArea(c *gin.Context) {
	var fromRequest secret_area.CreateRequest
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

	userId := c.GetInt("userId")
	if userId == 0 {
		c.Error(app_error.InternalServerError(errors.New("can't get user")))
		return
	}

	if data, err := h.Usecase.SecretArea.Create(fromRequest, userId); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// UpdateSecretArea godoc
// @Summary      Update secret area
// @Description  update secret area
// @Tags         secret-area
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	secret_area.UpdateRequest	true	"data"
// @Success      200	{object}	secret_area.SecretArea
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /secret-area/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) secretAreaUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest secret_area.UpdateRequest

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

	if data, err := h.Usecase.SecretArea.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// DeleteSecretArea godoc
// @Summary      Secret area delete
// @Description  secret area delete
// @Tags         secret-area
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.DeleteResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /secret-area/{id} [delete]
// @Security 	 BearerAuth
func (h *Handler) secretAreaDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.SecretArea.DeleteById(id); err != nil {
		c.Error(err)
	} else if data {
		c.JSON(http.StatusOK, transport.DeleteResponse{Success: true})
	} else {
		c.Error(app_error.ErrNotFound)
	}
}
