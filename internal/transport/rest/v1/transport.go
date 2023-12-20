package v1

import (
	"errors"
	"net/http"
	"seal/internal/domain/transport"
	transp "seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type transportList struct {
	RecordsTotal    int                   `json:"records_total"`
	RecordsFiltered int                   `json:"records_filtered"`
	Data            []transport.Transport `json:"data"`
}

func (h *Handler) registerTransportHandler(api *gin.RouterGroup) {
	group := api.Group("/transport")
	{

		group.GET(":id", h.transport)
		group.PUT(":id", h.transportUpdate)
		group.DELETE(":id", h.transportDelete)
		group.GET("", h.transportList)
		group.POST("", h.transportCreate)

	}
}

// ItemTransport godoc
// @Summary      Transport info
// @Description  transport info
// @Tags         transport
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.Transport
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /transport/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) transport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Transport.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListTransport godoc
// @Summary      List transports
// @Description  get transports
// @Tags         transport
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	transportList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /transport [get]
// @Security 	 BearerAuth
func (h *Handler) transportList(c *gin.Context) {
	var queryParams transp.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Transport.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// CreateTransport godoc
// @Summary      Create transport
// @Description  add transport
// @Tags         transport
// @Accept       json
// @Produce      json
// @Param		 data	body	transport.CreateRequest	true	"data"
// @Success      200	{object}	transport.Transport
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /transport [post]
// @Security 	 BearerAuth
func (h *Handler) transportCreate(c *gin.Context) {
	var fromRequest transport.CreateRequest
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

	if data, err := h.Usecase.Transport.Create(fromRequest, userId); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// UpdateTransport godoc
// @Summary      Update transport
// @Description  update transport
// @Tags         transport
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	transport.UpdateRequest	true	"data"
// @Success      200	{object}	transport.Transport
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /transport/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) transportUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest transport.UpdateRequest

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

	if data, err := h.Usecase.Transport.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// DeleteTransport godoc
// @Summary      Transport delete
// @Description  transport delete
// @Tags         transport
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.DeleteResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /transport/{id} [delete]
// @Security 	 BearerAuth
func (h *Handler) transportDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Transport.DeleteById(id); err != nil {
		c.Error(err)
	} else if data {
		c.JSON(http.StatusOK, transp.DeleteResponse{Success: true})
	} else {
		c.Error(app_error.ErrNotFound)
	}
}
