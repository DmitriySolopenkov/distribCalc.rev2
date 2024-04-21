package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code  int         `json:"code"`
	Error interface{} `json:"error"`
}

func Data(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &SuccessResponse{
		Code: http.StatusOK,
		Data: data,
	})
}

func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, &ErrorResponse{
		Code:  http.StatusBadRequest,
		Error: message,
	})
}

func InternalServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, &ErrorResponse{
		Code:  http.StatusInternalServerError,
		Error: message,
	})
}

func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, &ErrorResponse{
		Code:  http.StatusNotFound,
		Error: message,
	})
}
