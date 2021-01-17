package api

import (
	"github.com/gin-gonic/gin"
)

// @Deprecated
// @ID getMediaAttachments
// @Tags Attachments
// @Summary Get attachments of a media
// @Description Get attachments of a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Param media_id path string true "Media identifier" validate(required)
// @Router /medias/{media_id}/attachments [get]
func (s *Server) getMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}

// @Deprecated
// @ID addMediaAttachment
// @Tags Attachments
// @Summary Add attachment to a media
// @Description Add attachment to a media
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Param media_id path string true "Media identifier" validate(required)
// @Param key path string true "Attachment unique identifier" validate(required)
// @Param file formData file true "Attachment file"
// @Router /medias/{media_id}/attachments/{key} [post]
func (s *Server) addMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}

// @Deprecated
// @ID deleteMediaAttachment
// @Tags Attachments
// @Summary Delete attachment of a media
// @Description Delete attachment of a media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Param media_id path string true "Media identifier" validate(required)
// @Param key path string true "Attachment unique identifier" validate(required)
// @Router /medias/{media_id}/attachments/{key} [delete]
func (s *Server) deleteMediaAttachments(ctx *gin.Context) {
	ctx.Status(200)
}
