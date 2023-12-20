package middleware

import (
	"errors"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtWorker transport.JwtWorker) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := GetBearerToken(c.GetHeader("Authorization"))

		if err == nil && jwtWorker.TokenIsValid(token) {
			if userId, err := jwtWorker.GetUserIdFromToken(token); err != nil {
				c.Error(app_error.ErrUnauthorized)
				return
			} else {
				c.Set("userId", userId)
			}

			c.Next()
			return
		}

		c.Error(app_error.ErrUnauthorized)
		c.Abort()
	}
}

func GetBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Empty authorization header")
	}

	authHeaderParts := strings.Split(authHeader, "Bearer ")
	if len(authHeaderParts) != 2 {
		return "", errors.New("Wrong authorization header")
	}

	return authHeaderParts[1], nil
}
