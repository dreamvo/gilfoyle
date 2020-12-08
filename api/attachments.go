package api

import (
	"github.com/gin-gonic/gin"
)

// @ID getMediaAttachments
// @Tags Attachments
// @Summary Get attachments of a media
// @Description Get attachments of a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Param id path string true "Media identifier" validate(required)
// @Router /medias/{id}/attachments [get]
func getMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}

// @ID addMediaAttachment
// @Tags Attachments
// @Summary Add attachment to a media
// @Description Add attachment to a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Param id path string true "Media identifier" validate(required)
// @Param file formData file true "Attachment file"
// @Router /medias/{id}/attachments [post]
func addMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}

// @ID deleteMediaAttachment
// @Tags Attachments
// @Summary Delete attachment of a media
// @Description Delete attachment of a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Param id path string true "Media identifier" validate(required)
// @Param attachment_id path string true "Attachment identifier" validate(required)
// @Router /medias/{id}/attachments/{attachment_id} [delete]
func deleteMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}
