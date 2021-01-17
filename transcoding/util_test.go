package transcoding

import (
	"github.com/dreamvo/gilfoyle/ent"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtil(t *testing.T) {
	t.Run("ParseFrameRates", func(t *testing.T) {
		var suite = map[string]int8{
			"25/1":    25,
			"30/1":    30,
			"16/1":    16,
			"30":      30,
			"25":      25,
			"invalid": 0,
		}

		for input, expected := range suite {
			if fps := ParseFrameRates(input); fps != expected {
				t.Errorf("expected: %d, got: %d", expected, fps)
			}
		}
	})

	t.Run("CreateMasterPlaylist", func(t *testing.T) {
		t.Run("should", func(t *testing.T) {
			mediaFiles := []*ent.MediaFile{
				{
					RenditionName:    "720p",
					ResolutionWidth:  1280,
					ResolutionHeight: 720,
					TargetBandwidth:  800000,
				},
				{
					RenditionName:    "1080p",
					ResolutionWidth:  1920,
					ResolutionHeight: 1080,
					TargetBandwidth:  1000000,
				},
			}

			playlist := CreateMasterPlaylist(mediaFiles)

			assert.Equal(t, `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=800000,RESOLUTION=1280x720,NAME="720p"
720p/index.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=1000000,RESOLUTION=1920x1080,NAME="1080p"
1080p/index.m3u8
`, playlist)
		})
	})
}
