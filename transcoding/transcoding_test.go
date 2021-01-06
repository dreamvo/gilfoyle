package transcoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranscoding(t *testing.T) {
	t.Run("should return arguments", func(t *testing.T) {
		transcoder := NewTranscoder(Options{
			FFmpegBinPath:  "/usr/bin/ffmpeg",
			FFprobeBinPath: "/usr/bin/ffprobe",
		})

		format := "mp4"
		AudioCodec := "aac"
		VideoCodec := "h264"
		AudioRate := 48000
		Crf := uint32(20)
		KeyframeInterval := 48
		HlsSegmentDuration := 4
		HlsPlaylistType := "vod"
		VideoBitRate := 800000
		VideoMaxBitRate := 856000
		BufferSize := 1200000
		AudioBitrate := 96000
		SkipAudio := false

		p := transcoder.
			Process().
			Input("/tmp/input.mp4").
			Output("/tmp/output.mp4").
			WithOptions(ProcessOptions{
				OutputFormat:       &format,
				AudioCodec:         &AudioCodec,
				VideoCodec:         &VideoCodec,
				AudioRate:          &AudioRate,
				Crf:                &Crf,
				KeyframeInterval:   &KeyframeInterval,
				HlsSegmentDuration: &HlsSegmentDuration,
				HlsPlaylistType:    &HlsPlaylistType,
				VideoBitRate:       &VideoBitRate,
				VideoMaxBitRate:    &VideoMaxBitRate,
				BufferSize:         &BufferSize,
				AudioBitrate:       &AudioBitrate,
				SkipAudio:          &SkipAudio,
			}).
			WithAdditionalOptions(map[string]string{
				"key": "value",
			})

		assert.Equal(t, []string{"-b:v", "800000", "-maxrate", "856000", "-c:v", "h264", "-ar", "48000", "-g", "48", "-c:a", "aac", "-ab", "96000", "-bufsize", "1200000", "-f", "mp4", "-hls_playlist_type", "vod", "-hls_time", "4", "key", "value", "/tmp/output.mp4"}, p.GetStrArguments())
	})

	t.Run("should return arguments", func(t *testing.T) {
		transcoder := NewTranscoder(Options{
			FFmpegBinPath:  "/usr/bin/ffmpeg",
			FFprobeBinPath: "/usr/bin/ffprobe",
		})

		format := "mp4"
		AudioCodec := "aac"
		VideoCodec := "h264"
		AudioRate := 48000
		Crf := uint32(20)
		KeyframeInterval := 48
		HlsSegmentDuration := 4
		HlsPlaylistType := "vod"
		VideoBitRate := 800000
		VideoMaxBitRate := 856000
		BufferSize := 1200000
		AudioBitrate := 96000
		SkipAudio := true

		p := transcoder.
			Process().
			Input("/tmp/input.mp4").
			Output("/tmp/output.mp4").
			WithOptions(ProcessOptions{
				OutputFormat:       &format,
				AudioCodec:         &AudioCodec,
				VideoCodec:         &VideoCodec,
				AudioRate:          &AudioRate,
				Crf:                &Crf,
				KeyframeInterval:   &KeyframeInterval,
				HlsSegmentDuration: &HlsSegmentDuration,
				HlsPlaylistType:    &HlsPlaylistType,
				VideoBitRate:       &VideoBitRate,
				VideoMaxBitRate:    &VideoMaxBitRate,
				BufferSize:         &BufferSize,
				AudioBitrate:       &AudioBitrate,
				SkipAudio:          &SkipAudio,
			})

		assert.Equal(t, []string{"-b:v", "800000", "-maxrate", "856000", "-c:v", "h264", "-ar", "48000", "-g", "48", "-c:a", "aac", "-ab", "96000", "-bufsize", "1200000", "-f", "mp4", "-hls_playlist_type", "vod", "-hls_time", "4", "-an", "-i", "/tmp/input.mp4", "/tmp/output.mp4"}, p.GetStrArguments())
	})
}
