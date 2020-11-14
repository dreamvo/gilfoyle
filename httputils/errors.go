package httputils

import (
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"go.uber.org/zap"
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

// use single instances, it caches struct info
var (
	uni        *ut.UniversalTranslator
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	enTrans := en.New()
	uni = ut.New(enTrans, enTrans)

	trans, _ := uni.GetTranslator("en")

	translator = trans
	validate = validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}
}

// NewError returns a new error response
func NewError(ctx *gin.Context, status int, err error) {
	gilfoyle.Logger.Error("HTTP request resulted in an error", zap.Error(err))

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
		fields[strings.ToLower(err.Field())] = err.Translate(translator)
	}

	response := ValidationErrorResponse{
		Code:    ValidationErrorCode,
		Message: "Some parameters are missing or invalid",
		Fields:  fields,
	}
	ctx.JSON(ValidationErrorCode, response)
}
