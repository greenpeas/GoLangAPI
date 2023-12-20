package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seal/internal/domain/seal"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"
)

// List for swagger only
type listSeal struct {
	RecordsTotal    int `json:"records_total"`
	RecordsFiltered int `json:"records_filtered"`
	Data            []seal.SealForList
}

func (h *Handler) registerSealHandler(api *gin.RouterGroup) {
	group := api.Group("/seal")
	{
		group.GET(":id", h.seal)
		group.PUT(":id", h.sealUpdate)
		group.GET("", h.sealList)
		group.GET(":id/archive", h.sealArchive)
	}
}

// ItemSeal godoc
// @Summary      Seal info
// @Description  seal info
// @Tags         seal
// @Accept       json
// @Param        id       path     int     true  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	seal.Seal
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /seal/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) seal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Seal.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListSeal godoc
// @Summary      List seal
// @Description  get seals
// @Tags         seal
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	listSeal
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /seal [get]
// @Security 	 BearerAuth
func (h *Handler) sealList(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Seal.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// ArchiveSeal godoc
// @Summary      Archive seal
// @Description  get seal archive
// @Tags         seal
// @Accept       json
// @Param        id		     path    int     true  "id"		minimum(0)		maximum (32767)
// @Param        from	     query   string	 true  "from"
// @Param        to	    	 query	 string	 false "to"
// @Param        limit		 query   int   	 false "limit"
// @Param        order_desc  query	 bool	 false "order_desc"
// @Success      200	{object}	[]sealData.SealData
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /seal/{id}/archive [get]
// @Security 	 BearerAuth
func (h *Handler) sealArchive(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	queryParams := seal.ArchiveQueryParams{Id: id}
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if data, err := h.Usecase.Seal.Archive(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// UpdateSeal godoc
// @Summary      Update seal
// @Description  update seal
// @Tags         seal
// @Accept       json
// @Produce      json
// @Param        id     path        int     true  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	    seal.UpdateRequest	true	"data"
// @Success      200	{object}	seal.Seal
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /seal/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) sealUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest seal.UpdateRequest

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

	if createdRoute, err := h.Usecase.Seal.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, createdRoute)
	}
}
