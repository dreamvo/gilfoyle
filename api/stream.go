package api

import (
	"github.com/gin-gonic/gin"
)

// @ID streamMedia
// @Tags Stream
// @Summary Get stream from media file
// @Description Get stream from media file
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Header 200 {string} Content-Type "application/octet-stream"
// @Param media_id path string true "Media identifier" validate(required)
// @Param preset path string true "Encoder preset" validate(required)
// @Router /medias/{media_id}/stream/{preset} [get]
func streamMedia(ctx *gin.Context) {
	ctx.Status(200)
}
