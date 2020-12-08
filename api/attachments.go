package api

import (
	"github.com/gin-gonic/gin"
)

type AddMediaAttachment struct {
	Key string `json:"key" validate:"required,gte=1,lte=100" example:"subtitle_fr_FR"`
}

// @ID getMediaAttachments
// @Tags Attachments
// @Summary Get attachments of a media
// @Description Get attachments of a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Param media_id path string true "Media identifier" validate(required)
// @Router /medias/{media_id}/attachments [get]
func getMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}

// @ID addMediaAttachment
// @Tags Attachments
// @Summary Add attachment to a media
// @Description Add attachment to a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Param media_id path string true "Media identifier" validate(required)
// @Param file formData file true "Attachment file"
// @Param attachment body AddMediaAttachment true "Attachment metadata" validate(required)
// @Router /medias/{media_id}/attachments [post]
func addMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}

// @ID deleteMediaAttachment
// @Tags Attachments
// @Summary Delete attachment of a media
// @Description Delete attachment of a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Param media_id path string true "Media identifier" validate(required)
// @Param attachment_id path string true "Attachment identifier" validate(required)
// @Router /medias/{media_id}/attachments/{attachment_id} [delete]
func deleteMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}
