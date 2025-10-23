package util

import (
	"net/http"

	"user-service/internal/model"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, model.ErrorResponse{
		Status:  "error",
		Message: message,
		Data:    errors,
	})
}

func RespondWithSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, model.APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}
