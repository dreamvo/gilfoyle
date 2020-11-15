package util

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidateBody(ctx *gin.Context, obj interface{}) error {
	err := ctx.BindJSON(&obj)
	if err != nil {
		return err
	}

	if err := validate.Struct(obj); err != nil {
		return err
	}

	return nil
}

func ValidateUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	return parsedUUID, err
}
