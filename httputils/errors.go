package httputils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse example
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type ValidationField struct {
	Message string
	Type    string
}

type ValidationErrorResponse struct {
	Code    int                        `json:"code" example:"400"`
	Message string                     `json:"message" example:"status bad request"`
	Fields  map[string]ValidationField `json:"fields"`
}

// NewError returns a new error response
func NewError(ctx *gin.Context, status int, err error) {
	er := ErrorResponse{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// NewValidationError returns a new validation error response
func NewValidationError(ctx *gin.Context, status int, err error) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		NewError(ctx, 500, fmt.Errorf("Unexpected error occurred"))
		return
	}

	fields := map[string]ValidationField{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = ValidationField{
			Message: "field is invalid",
			Type:    err.Type().String(),
		}
	}

	response := ValidationErrorResponse{
		Code:    status,
		Message: "Some parameters are missing or invalid",
		Fields:  fields,
	}
	ctx.JSON(status, response)
}
