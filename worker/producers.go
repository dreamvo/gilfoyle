package worker

import (
	"encoding/json"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type VideoTranscodingParams struct {
	OriginalFile       transcoding.OriginalFile `json:"original_file"`
	MediaUUID          uuid.UUID                `json:"media_uuid"`
	RenditionName      string                   `json:"preset_name"`
	VideoWidth         int                      `json:"video_width"`
	VideoHeight        int                      `json:"video_height"`
	AudioCodec         string                   `json:"audio_codec"`
	VideoCodec         string                   `json:"video_codec"`
	Crf                uint32                   `json:"crf"`
	KeyframeInterval   int                      `json:"keyframe_interval"`
	HlsSegmentDuration int                      `json:"hls_segment_duration"`
	HlsPlaylistType    string                   `json:"hls_playlist_type"`
	VideoBitRate       int                      `json:"video_bit_rate"`
	AudioBitrate       int                      `json:"audio_bitrate"`
	FrameRate          int                      `json:"frame_rate"`
	TargetBandwidth    uint64                   `json:"target_bandwidth"`
}

type MediaProcessingCallbackParams struct {
	MediaUUID       uuid.UUID `json:"media_uuid"`
	MediaFilesCount int       `json:"media_files_count"`
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

func MediaProcessingCallbackProducer(ch Channel, data MediaProcessingCallbackParams) error {
	body, _ := json.Marshal(data)

	err := ch.Publish("", MediaProcessingCallbackQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
