package v1

import (
	"net/http"
	"seal/internal/domain/route"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type routeList struct {
	RecordsTotal    int           `json:"records_total"`
	RecordsFiltered int           `json:"records_filtered"`
	Data            []route.Route `json:"data"`
}

func (h *Handler) registerRouteHandler(api *gin.RouterGroup) {
	group := api.Group("/route")
	{
		group.GET(":id", h.route)
		group.PUT(":id", h.routeUpdate)
		group.DELETE(":id", h.routeDelete)
		group.GET("", h.routeList)
		group.POST("", h.routeCreate)
	}
}

// ItemRoute godoc
// @Summary      Route info
// @Description  route info
// @Tags         route
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	route.Route
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /route/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) route(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Route.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListRoutes godoc
// @Summary      List router
// @Description  get routes
// @Tags         route
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	routeList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /route [get]
// @Security 	 BearerAuth
func (h *Handler) routeList(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Route.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// CreateRoute godoc
// @Summary      Create route
// @Description  add route
// @Tags         route
// @Accept       json
// @Produce      json
// @Param		 data	body	route.CreateRequest	true	"data"
// @Success      200	{object}	route.Route
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /route [post]
// @Security 	 BearerAuth
func (h *Handler) routeCreate(c *gin.Context) {
	var fromRequest route.CreateRequest
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

	if createdRoute, err := h.Usecase.Route.Create(fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdRoute)
	}
}

// UpdateRoute godoc
// @Summary      Update route
// @Description  update route
// @Tags         route
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	route.UpdateRequest	true	"data"
// @Success      200	{object}	route.Route
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /route/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) routeUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest route.UpdateRequest

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

	if createdRoute, err := h.Usecase.Route.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdRoute)
	}
}

// ItemRoute godoc
// @Summary      Route delete
// @Description  route delete
// @Tags         route
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.DeleteResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /route/{id} [delete]
// @Security 	 BearerAuth
func (h *Handler) routeDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Route.DeleteById(id); err != nil {
		c.Error(err)
	} else if data {
		c.JSON(http.StatusOK, transport.DeleteResponse{Success: true})
	} else {
		c.Error(app_error.ErrNotFound)
	}
}
