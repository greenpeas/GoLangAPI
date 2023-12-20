package v1

import (
	"net/http"
	"seal/internal/domain/custom"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type customList struct {
	RecordsTotal    int             `json:"records_total"`
	RecordsFiltered int             `json:"records_filtered"`
	Data            []custom.Custom `json:"data"`
}

func (h *Handler) registerCustomHandler(api *gin.RouterGroup) {
	group := api.Group("/custom")
	{
		group.GET(":id", h.custom)
		group.PUT(":id", h.customUpdate)
		group.DELETE(":id", h.customDelete)
		group.GET("", h.customList)
		group.POST("", h.customCreate)
	}
}

// ItemCustom godoc
// @Summary      Custom info
// @Description  custom info
// @Tags         custom
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	custom.Custom
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /custom/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) custom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Custom.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListCustom godoc
// @Summary      List custom
// @Description  get customs
// @Tags         custom
// @Accept       json
// @Produce      json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	customList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /custom [get]
// @Security 	 BearerAuth
func (h *Handler) customList(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Custom.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// CreateCustom godoc
// @Summary      Create custom
// @Description  add custom
// @Tags         custom
// @Accept       json
// @Produce      json
// @Param		 data	body	custom.CreateRequest	true	"data"
// @Success      200	{object}	custom.Custom
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /custom [post]
// @Security 	 BearerAuth
func (h *Handler) customCreate(c *gin.Context) {
	var fromRequest custom.CreateRequest
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

	if createdRoute, err := h.Usecase.Custom.Create(fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdRoute)
	}
}

// UpdateCustom godoc
// @Summary      Update custom
// @Description  update custom
// @Tags         custom
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	custom.UpdateRequest	true	"data"
// @Success      200	{object}	custom.Custom
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /custom/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) customUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest custom.UpdateRequest

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

	if createdCustom, err := h.Usecase.Custom.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdCustom)
	}
}

// ItemRoute godoc
// @Summary      Custom delete
// @Description  Custom delete
// @Tags         custom
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.DeleteResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /custom/{id} [delete]
// @Security 	 BearerAuth
func (h *Handler) customDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Custom.DeleteById(id); err != nil {
		c.Error(err)
	} else if data {
		c.JSON(http.StatusOK, transport.DeleteResponse{Success: true})
	} else {
		c.Error(app_error.ErrNotFound)
	}
}
