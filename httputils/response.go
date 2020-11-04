package httputils

import "github.com/gin-gonic/gin"

// DataResponse example
type DataResponse struct {
	Code int         `json:"code" example:"200"`
	Data interface{} `json:"data,omitempty"`
}

// NewData returns a new response following the DataResponse schema
func NewData(ctx *gin.Context, status int, data interface{}) {
	res := DataResponse{
		Code: status,
		Data: data,
	}
	ctx.JSON(status, res)
}

// NewResponse returns a new response
func NewResponse(ctx *gin.Context, status int, obj interface{}) {
	ctx.JSON(status, obj)
}
