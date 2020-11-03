package httputils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

const (
	ValidationErrorCode int = 400
)

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// ValidationErrorResponse represents a validation error API response
type ValidationErrorResponse struct {
	Code    int               `json:"code" example:"400"`
	Message string            `json:"message" example:"status bad request"`
	Fields  map[string]string `json:"fields"`
}

// NewError returns a new error response
func NewError(ctx *gin.Context, status int, err error) {
	response := ErrorResponse{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, response)
}

// NewValidationError returns a new validation error response
func NewValidationError(ctx *gin.Context, err error) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		NewError(ctx, ValidationErrorCode, fmt.Errorf("Bad request"))
		return
	}

	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[strings.ToLower(err.Field())] = err.Error()
	}

	response := ValidationErrorResponse{
		Code:    ValidationErrorCode,
		Message: "Some parameters are missing or invalid",
		Fields:  fields,
	}
	ctx.JSON(ValidationErrorCode, response)
}
