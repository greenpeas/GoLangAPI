package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"seal/internal/domain/modem"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"
	"time"
)

type modemSwagger struct {
	Id          int          `json:"id,omitempty"`
	Imei        string       `json:"imei,omitempty"`
	Serial      int          `json:"serial,omitempty"`
	Iccid       string       `json:"iccid,omitempty"`
	LastDevTime time.Time    `json:"last_dev_time"`
	Extra       *modem.Extra `json:"extra,omitempty"`
	Last        *modem.Data  `json:"last,omitempty"`
	Seals       []modem.Seal `json:"seals,omitempty"`
}

// List for swagger only
type listModem struct {
	RecordsTotal    int `json:"records_total"`
	RecordsFiltered int `json:"records_filtered"`
	Data            []modem.ModemForList
}
type listModemShippingReady struct {
	RecordsTotal    int `json:"records_total"`
	RecordsFiltered int `json:"records_filtered"`
	Data            []modem.ModemForListShippingReady
}

func (h *Handler) registerModemHandler(api *gin.RouterGroup) {
	group := api.Group("/modem")
	{
		group.GET(":id", h.modem)
		group.PUT(":id", h.modemUpdate)
		group.GET("", h.modemList)
		group.GET("shipping-ready", h.modemListShippingReady)
		group.POST(":id/command", h.modemSendCommand)
		group.GET(":id/commands", h.modemListCommands)
		group.GET(":id/archive", h.modemArchive)
		group.GET(":id/log-raw-telemetry", h.modemLogRawTelemetry)
		group.GET(":id/track", h.modemTrack)
		group.GET(":id/log", h.modemLog)
	}
}

// ItemModem godoc
// @Summary      Modem info (by id or imei)
// @Description  modem info (by id or imei)
// @Tags         modem
// @Accept       json
// @Param        id       path     int     false  "id"
// @Success      200	{object}	modemSwagger
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) modem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if id > math.MaxInt32 {
		if data, err := h.Usecase.Modem.GetByImei(id); err != nil {
			c.Error(err)
		} else {
			c.JSON(http.StatusOK, data)
		}

		return
	}

	if data, err := h.Usecase.Modem.GetById(int(id)); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}

}

// ListModem godoc
// @Summary      List modem
// @Description  get modems
// @Tags         modem
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	listModem
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem [get]
// @Security 	 BearerAuth
func (h *Handler) modemList(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Modem.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// ListModemReadyForShipping godoc
// @Summary      List modem ready for shipping
// @Description  get modems ready for shipping
// @Tags         modem
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	listModemShippingReady
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/shipping-ready [get]
// @Security 	 BearerAuth
func (h *Handler) modemListShippingReady(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Modem.ListShippingReady(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// SendCommandModem godoc
// @Summary      Send command modem
// @Description  send command for modem
// @Tags         modem
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	modem.SendCommandRequest	true	"data"
// @Success      200	{object}	transport.SuccessResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id}/command [post]
// @Security 	 BearerAuth
func (h *Handler) modemSendCommand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest modem.SendCommandRequest
	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	userId := c.GetInt("userId")
	if userId == 0 {
		c.Error(app_error.InternalServerError(errors.New("can't get user")))
		return
	}
	user, uerr := h.Usecase.User.GetById(userId)
	if uerr != nil {
		c.Error(app_error.InternalServerError(errors.New("can't get user by id")))
		return
	}

	if data, err := h.Usecase.Modem.SendCommand(id, fromRequest, user.Login); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, transport.SuccessResponse{Success: data})
	}
}

// ListModemCommands godoc
// @Summary      List modem commands
// @Description  get modem commands
// @Tags         modem
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	listModem
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id}/commands [get]
// @Security 	 BearerAuth
func (h *Handler) modemListCommands(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if list, err := h.Usecase.Modem.CommandsList(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// ArchiveModem godoc
// @Summary      Archive modem
// @Description  get modem archive
// @Tags         modem
// @Accept       json
// @Param        id			 path    int     true  "id"		minimum(0)		maximum (32767)
// @Param        from		 query	 string	 true  "from"
// @Param        to	    	 query	 string	 false "to"
// @Param        limit		 query   int   	 false "limit"
// @Param        order_desc  query	 bool	 false "order_desc"
// @Success      200	{object}	[]modem.ArchiveModemData
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id}/archive [get]
// @Security 	 BearerAuth
func (h *Handler) modemArchive(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	queryParams := modem.ArchiveQueryParams{Id: id}
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if data, err := h.Usecase.Modem.Archive(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ArchiveModemByLogRaw godoc
// @Summary      Archive modem from raw log
// @Description  get modem archive
// @Tags         modem
// @Accept       json
// @Param        id			 path    int     true  "id"		minimum(0)		maximum (32767)
// @Param        from		 query	 string	 true  "from"
// @Param        to	    	 query	 string	 false "to"
// @Param        limit		 query   int   	 false "limit"
// @Param        order_desc  query	 bool	 false "order_desc"
// @Success      200	{object}	[]modem.ArchiveModemData
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id}/log-raw-telemetry [get]
// @Security 	 BearerAuth
func (h *Handler) modemLogRawTelemetry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	queryParams := modem.ArchiveQueryParams{Id: id}
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if data, err := h.Usecase.Modem.LogRawTelemetry(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// TrackModem godoc
// @Summary      Track modem
// @Description  get modem coordinates track
// @Tags         modem
// @Accept       json
// @Param        id			 path    int    true  "id"		minimum(0)		maximum (32767)
// @Param        from		 query	 string	true  "from"
// @Param        to	    	 query	 string	false "to"
// @Param        limit		 query   int   	false "limit"
// @Param        order_desc  query	 bool	false "order_desc"
// @Success      200	{object}	modem.TrackResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id}/track [get]
// @Security 	 BearerAuth
func (h *Handler) modemTrack(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	queryParams := modem.TrackQueryParams{Id: id}
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if data, err := h.Usecase.Modem.Track(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// LogModem godoc
// @Summary      Log modem
// @Description  get modem packages log
// @Tags         modem
// @Accept       json
// @Param        id		     path    int     true  "id"		minimum(0)		maximum (32767)
// @Param        from	     query	 string	 true  "from"
// @Param        to	    	 query	 string	 false "to"
// @Param        limit		 query   int   	 false "limit"
// @Param        order_desc  query	 bool	 false "order_desc"
// @Success      200	{object}	[]modemLogRaw.ModemLogRaw
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id}/log [get]
// @Security 	 BearerAuth
func (h *Handler) modemLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	queryParams := modem.LogQueryParams{Id: id}
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if data, err := h.Usecase.Modem.Log(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// UpdateModem godoc
// @Summary      Update modem
// @Description  update modem
// @Tags         modem
// @Accept       json
// @Produce      json
// @Param        id     path        int     true  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	    modem.UpdateRequest	true	"data"
// @Success      200	{object}	modem.Modem
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /modem/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) modemUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest modem.UpdateRequest

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

	if createdRoute, err := h.Usecase.Modem.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdRoute)
	}
}
