package api

import (
	_ "github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/httputils"
	"github.com/gin-gonic/gin"
)

// @ID uploadMediaFile
// @Tags Medias
// @Summary Upload a media file
// @Description Upload a new media file for a given media ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=ent.Media}
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /medias/{id}/upload [post]
// @Param id path string true "Media identifier" validate(required)
// @Param file formData file true "Media file"
func uploadMediaFile(ctx *gin.Context) {
	ctx.Status(200)
}
