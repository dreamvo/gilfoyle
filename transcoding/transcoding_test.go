package transcoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranscoding(t *testing.T) {
	t.Run("should return arguments (1)", func(t *testing.T) {
		transcoder := NewTranscoder(Options{
			FFmpegBinPath: "/usr/bin/ffmpeg",
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
			SetInput("/tmp/input.mp4").
			SetOutput("/tmp/output.mp4").
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
				WhiteListProtocols: []string{"file", "http", "https", "tcp", "tls"},
			}).
			WithAdditionalOptions(map[string]string{
				"key": "value",
			})

		assert.Equal(t, []string{"-i", "/tmp/input.mp4", "-b:v", "800000", "-maxrate", "856000", "-c:v", "h264", "-ar", "48000", "-g", "48", "-c:a", "aac", "-ab", "96000", "-bufsize", "1200000", "-f", "mp4", "-hls_playlist_type", "vod", "-hls_time", "4", "-protocol_whitelist", "file,http,https,tcp,tls", "key", "value", "/tmp/output.mp4"}, p.GetStrArguments())
	})

	t.Run("should return arguments (2)", func(t *testing.T) {
		transcoder := NewTranscoder(Options{
			FFmpegBinPath: "/usr/bin/ffmpeg",
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
			SetInput("/tmp/input.mp4").
			SetOutput("/tmp/output.mp4").
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

		assert.Equal(t, []string{"-i", "/tmp/input.mp4", "-b:v", "800000", "-maxrate", "856000", "-c:v", "h264", "-ar", "48000", "-g", "48", "-c:a", "aac", "-ab", "96000", "-bufsize", "1200000", "-f", "mp4", "-hls_playlist_type", "vod", "-hls_time", "4", "-an", "/tmp/output.mp4"}, p.GetStrArguments())
	})

	t.Run("should return arguments (3)", func(t *testing.T) {
		transcoder := NewTranscoder(Options{
			FFmpegBinPath: "/usr/bin/ffmpeg",
		})

		p := transcoder.
			Process().
			SetInput("/tmp/input.mp4").
			SetOutput("/tmp/output.ts").
			WithOptions(ProcessOptions{
				StreamIds: map[string]string{
					"0": "33",
					"1": "36",
				},
			})

		assert.Contains(t, [][]string{
			{"-i", "/tmp/input.mp4", "-streamid", "0:33", "-streamid", "1:36", "/tmp/output.ts"},
			{"-i", "/tmp/input.mp4", "-streamid", "1:36", "-streamid", "1:33", "/tmp/output.ts"},
		}, p.GetStrArguments())
	})

	t.Run("should fail to execute command", func(t *testing.T) {
		transcoder := NewTranscoder(Options{
			FFmpegBinPath: "/test/ffmpeg",
		})

		p := transcoder.
			Process().
			SetInput("/tmp/input.mp4").
			SetOutput("/tmp/output.ts")

		err := transcoder.Run(p)
		assert.Error(t, err)
		assert.EqualError(t, err, "fork/exec /test/ffmpeg: no such file or directory")
	})
}
