package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @ID getMediaMasterPlaylist
// @Tags Stream
// @Summary Get HLS master playlist of media
// @Description Get HLS master playlist of media
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Header 200 {string} Content-Type "application/octet-stream"
// @Param media_id path string true "Media identifier" validate(required)
// @Router /medias/{media_id}/stream/playlist [get]
func (s *Server) getMediaMasterPlaylist(ctx *gin.Context) {
	mediaUUID := ctx.Param("id")

	parsedUUID, err := util.ValidateUUID(mediaUUID)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, ErrInvalidUUID)
		return
	}

	m, _ := s.db.Media.Query().Where(media.ID(parsedUUID)).WithMediaFiles().Only(context.Background())
	if m == nil {
		util.NewError(ctx, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	if m.Status != schema.MediaStatusReady {
		util.NewError(ctx, http.StatusTooEarly, errors.New("media is not ready yet for streaming"))
		return
	}

	masterPlaylistPath := fmt.Sprintf("%s/%s", m.ID.String(), transcoding.HLSMasterPlaylistFilename)

	// Create master playlist if it does not exists
	stat, _ := s.storage.Stat(context.Background(), masterPlaylistPath)
	if stat == nil {
		err = s.storage.Save(context.Background(), strings.NewReader(transcoding.CreateMasterPlaylist(m.Edges.MediaFiles)), masterPlaylistPath)
		if err != nil {
			util.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Open the master playlist file
	f, err := s.storage.Open(context.Background(), masterPlaylistPath)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	stat, err = s.storage.Stat(context.Background(), masterPlaylistPath)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.DataFromReader(http.StatusOK, stat.Size, "application/octet-stream", f, map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, transcoding.HLSMasterPlaylistFilename),
	})
}

// @ID getMediaPlaylist
// @Tags Stream
// @Summary Get file of HLS playlist
// @Description Get file of HLS playlist
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Header 200 {string} Content-Type "application/octet-stream"
// @Param media_id path string true "Media identifier" validate(required)
// @Param playlist path string true "HLS playlist name" validate(required)
// @Param filename path string true "HLS filename" validate(required)
// @Router /medias/{media_id}/stream/playlist/{playlist}/{filename} [get]
func (s *Server) getMediaPlaylist(ctx *gin.Context) {
	mediaUUID := ctx.Param("id")
	playlistName := ctx.Param("playlist")
	filename := ctx.Param("filename")

	parsedUUID, err := util.ValidateUUID(mediaUUID)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, ErrInvalidUUID)
		return
	}

	m, _ := s.db.Media.Query().Where(media.ID(parsedUUID)).WithMediaFiles().Only(context.Background())
	if m == nil {
		util.NewError(ctx, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	if m.Status != schema.MediaStatusReady {
		util.NewError(ctx, http.StatusTooEarly, errors.New("media is not ready yet for streaming"))
		return
	}

	f, err := s.storage.Open(context.Background(), fmt.Sprintf("%s/%s/%s", m.ID.String(), playlistName, filename))
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	stat, err := s.storage.Stat(context.Background(), fmt.Sprintf("%s/%s/%s", m.ID.String(), playlistName, filename))
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.DataFromReader(http.StatusOK, stat.Size, "application/octet-stream", f, map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, filename),
	})
}
