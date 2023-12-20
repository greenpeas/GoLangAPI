package middleware

import (
	"net/http"
	app_interface "seal/internal/app/interface"
	"seal/pkg/app_error"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()

		if err == nil {
			return
		}

		appError, ok := err.Unwrap().(app_interface.Error)

		if !ok || appError == nil {
			c.JSON(http.StatusInternalServerError, app_error.InternalServerError(err))
			return
		}

		//if appError.Unwrap() != nil {
		//	logger.Error(appError.Unwrap().Error())
		//}

		if _, ok := appError.GetBody().(string); ok {
			c.JSON(appError.GetHttpCode(), appError)
		} else {
			c.JSON(appError.GetHttpCode(), appError.GetBody())
		}

	}
}
