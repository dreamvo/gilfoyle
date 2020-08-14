package httputils

import "github.com/gin-gonic/gin"

// ErrorResponse example
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// NewError returns a new error response
func NewError(ctx *gin.Context, status int, err error) {
	er := ErrorResponse{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}
