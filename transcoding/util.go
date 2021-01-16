package transcoding

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/ent"
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
	content := "#EXTM3U\n#EXT-X-VERSION:3\n"

	for _, mediaFile := range mediaFiles {
		content += fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s/%s\n",
			mediaFile.TargetBandwidth,
			mediaFile.ResolutionWidth,
			mediaFile.ResolutionHeight,
			mediaFile.RenditionName,
			HLSPlaylistFilename,
		)
	}

	return content
}
