package api

import (
	"context"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/x/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"testing"
)

func TestStream(t *testing.T) {
	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() { _ = dbClient.Close() }()

	_, err := gilfoyle.NewConfig()
	if err != nil {
		t.Error(err)
	}

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer removeDir(gilfoyle.Config.Storage.Filesystem.DataPath)

	storageDriver, err := gilfoyle.NewStorage(storage.Filesystem)
	if err != nil {
		t.Error(err)
	}

	s := NewServer(Options{
		Logger:   zap.NewExample(),
		Storage:  storageDriver,
		Database: dbClient,
	})

	t.Run("GET /medias/{media_id}/stream/playlist", func(t *testing.T) {
		t.Run("should create then return master playlist", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusReady).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("720p").
				SetResolutionWidth(1280).
				SetResolutionHeight(720).
				SetFormat("hls").
				SetTargetBandwidth(2800000).
				SetVideoBitrate(800000).
				SetFramerate(30).
				SetDurationSeconds(5).
				SetMediaType(schema.MediaFileTypeVideo).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("1080p").
				SetResolutionWidth(1920).
				SetResolutionHeight(1080).
				SetFormat("hls").
				SetTargetBandwidth(5000000).
				SetVideoBitrate(800000).
				SetFramerate(30).
				SetDurationSeconds(5).
				SetMediaType(schema.MediaFileTypeVideo).
				Save(context.Background())
			assert.NoError(t, err)

			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/playlist", m.ID.String()), nil)
			assert.NoError(t, err)

			stat, err := storageDriver.Stat(context.Background(), fmt.Sprintf("%s/%s", m.ID.String(), transcoding.HLSMasterPlaylistFilename))
			assert.NoError(t, err)

			assert.Equal(t, int64(171), stat.Size)

			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
			assert.Equal(t, `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1280x720
720p/index.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080
1080p/index.m3u8
`, res.Body.String())
		})

		t.Run("should return playlist entry file", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusReady).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("720p").
				SetResolutionWidth(1280).
				SetResolutionHeight(720).
				SetFormat("hls").
				SetTargetBandwidth(2800000).
				SetVideoBitrate(800000).
				SetFramerate(30).
				SetDurationSeconds(5).
				SetMediaType(schema.MediaFileTypeVideo).
				Save(context.Background())
			assert.NoError(t, err)

			playlistContent := `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:5
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-PLAYLIST-TYPE:VOD
#EXTINF:5.280000,
000.ts
#EXT-X-ENDLIST
`

			err = storageDriver.Save(context.Background(), strings.NewReader(playlistContent), fmt.Sprintf("%s/%s/%s", m.ID.String(), "360p", "index.m3u8"))
			assert.NoError(t, err)

			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/playlist/360p/index.m3u8", m.ID.String()), nil)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
			assert.Equal(t, playlistContent, res.Body.String())
		})
	})
}
