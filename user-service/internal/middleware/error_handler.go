package middleware

import (
	"net/http"

	"user-service/internal/util"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			util.RespondWithError(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
}
