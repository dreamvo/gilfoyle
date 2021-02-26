package util

import "github.com/gin-gonic/gin"

// DataResponse example
type DataResponse struct {
	Code     int         `json:"code" example:"200"`
	Metadata interface{} `json:"metadata,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// NewData returns a new response following the DataResponse schema
func NewData(ctx *gin.Context, status int, data interface{}, metadata interface{}) {
	res := DataResponse{
		Code:     status,
		Metadata: metadata,
		Data:     data,
	}
	ctx.JSON(status, res)
}

// NewResponse returns a new response
func NewResponse(ctx *gin.Context, status int, obj interface{}) {
	ctx.JSON(status, obj)
}
