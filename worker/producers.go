package worker

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type VideoTranscodingParams struct {
	MediaUUID          uuid.UUID `json:"media_uuid"`
	OriginalFilePath   string    `json:"source_file_path"`
	RenditionName      string    `json:"preset_name"`
	VideoWidth         int       `json:"video_width"`
	VideoHeight        int       `json:"video_height"`
	AudioCodec         string    `json:"audio_codec"`
	AudioRate          int       `json:"audio_rate"`
	VideoCodec         string    `json:"video_codec"`
	Crf                uint32    `json:"crf"`
	KeyframeInterval   int       `json:"keyframe_interval"`
	HlsSegmentDuration int       `json:"hls_segment_duration"`
	HlsPlaylistType    string    `json:"hls_playlist_type"`
	VideoBitRate       int       `json:"video_bit_rate"`
	VideoMaxBitRate    int       `json:"video_max_bit_rate"`
	BufferSize         int       `json:"buffer_size"`
	AudioBitrate       int       `json:"audio_bitrate"`
	FrameRate          int       `json:"frame_rate"`
}

func VideoTranscodingProducer(ch Channel, data VideoTranscodingParams) error {
	body, _ := json.Marshal(data)

	err := ch.Publish("", VideoTranscodingQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
