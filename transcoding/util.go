package transcoding

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/ent"
	"github.com/grafov/m3u8"
	"path"
	"strconv"
	"strings"
)

func ParseFrameRates(f string) int8 {
	slice := strings.Split(f, "/")
	i, err := strconv.ParseInt(slice[0], 10, 16)
	if err != nil {
		return 0
	}

	return int8(i)
}

func CreateMasterPlaylist(mediaFiles []*ent.MediaFile) string {
	p := m3u8.NewMasterPlaylist()

	for _, mediaFile := range mediaFiles {
		p.Append(
			path.Join(mediaFile.RenditionName, HLSPlaylistFilename),
			&m3u8.MediaPlaylist{
				TargetDuration: mediaFile.DurationSeconds,
			},
			m3u8.VariantParams{
				Name:       mediaFile.RenditionName,
				Bandwidth:  uint32(mediaFile.TargetBandwidth),
				FrameRate:  float64(mediaFile.Framerate),
				Resolution: fmt.Sprintf("%dx%d", mediaFile.ResolutionWidth, mediaFile.ResolutionHeight),
			},
		)
	}

	return p.Encode().String()
}
