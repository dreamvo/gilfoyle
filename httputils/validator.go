package httputils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ValidateBody(ctx *gin.Context, obj interface{}) error {
	err := ctx.BindJSON(&obj)
	if err != nil {
		return err
	}

	if err := validator.New().Struct(obj); err != nil {
		return err
	}

	return nil
}

func ValidateUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	return parsedUUID, err
}
