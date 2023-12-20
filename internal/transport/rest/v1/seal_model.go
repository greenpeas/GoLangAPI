package v1

import (
	"net/http"
	"seal/internal/domain/seal_model"
	"seal/internal/transport"
	"seal/pkg/app_error"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type sealModelList struct {
	RecordsTotal    int                    `json:"records_total"`
	RecordsFiltered int                    `json:"records_filtered"`
	Data            []seal_model.SealModel `json:"data"`
}

func (h *Handler) registerSealModelHandler(api *gin.RouterGroup) {
	group := api.Group("/seal-model")
	{
		group.GET("", h.sealModelList)
	}
}

// ListSealModel godoc
// @Summary      List seal models
// @Description  get seal models
// @Tags         seal-model
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Success      200	{object}	sealModelList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /seal-model [get]
// @Security 	 BearerAuth
func (h *Handler) sealModelList(c *gin.Context) {
	var queryParams transport.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.SealModel.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}
