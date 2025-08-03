package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		StatusCode: http.StatusCreated,
		Message:    message,
		Data:       data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       map[string]interface{}{},
	})
}

func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message)
}

func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}