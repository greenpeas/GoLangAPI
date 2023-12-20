package middleware

import (
	"bytes"
	"fmt"
	"io"
	app_interface "seal/internal/app/interface"

	"github.com/gin-gonic/gin"
)

type mes struct {
	method string
	url    string
	header any
	body   any
}

func Logger(logger app_interface.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)

		message := fmt.Sprintf("%+v", mes{
			method: c.Request.Method,
			url:    c.Request.RequestURI,
			header: c.Request.Header,
			body:   string(body),
		})

		logger.Debug(message)

		c.Next()
	}
}
