package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) registerShippingTempHandler(api *gin.RouterGroup) {
	group := api.Group("/shipping-temp")
	{
		group.POST(":id/files/:type", h.shippingUploadFiles)
	}
}

// UploadShippingFiles godoc
// @Summary      Upload shipping files
// @Description  upload shipping files
// @Tags         shipping-temp
// @Accept       multipart/form-data
// @Param        id       path     int     true  "id"		minimum(0)	maximum (32767)
// @Param        type     path     int     true  "type" 	minimum(0)	maximum (3)
// @Param        file   formData    []file true  "upload files"
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping-temp/{id}/files/{type} [post]
// @Security 	 BearerAuth
func (h *Handler) f1(c *gin.Context) {

}
