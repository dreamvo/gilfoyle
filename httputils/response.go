package httputils

import "github.com/gin-gonic/gin"

// HTTPResponse example
type HTTPResponse struct {
	Code int         `json:"code" example:"200"`
	Data interface{} `json:"data,omitempty"`
}

// NewResponse returns a new response
func NewResponse(ctx *gin.Context, status int, data interface{}) {
	er := HTTPResponse{
		Code: status,
		Data: data,
	}
	ctx.JSON(status, er)
}
