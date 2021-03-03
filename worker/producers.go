package worker

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type HlsVideoEncodingParams struct {
	MediaFileUUID      uuid.UUID `json:"media_file_uuid"`
	KeyframeInterval   int       `json:"keyframe_interval"`
	HlsSegmentDuration int       `json:"hls_segment_duration"`
	HlsPlaylistType    string    `json:"hls_playlist_type"`
}

type MediaEncodingCallbackParams struct {
	MediaUUID       uuid.UUID `json:"media_uuid"`
	MediaFilesCount int       `json:"media_files_count"`
}

type MediaEncodingEntrypoint struct {
	MediaUUID uuid.UUID `json:"media_uuid"`
}

func MediaEncodingEntrypointProducer(ch Channel, data MediaEncodingEntrypoint) error {
	body, _ := json.Marshal(data)

	err := ch.Publish("", EncodingEntrypointQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}

func HlsVideoEncodingProducer(ch Channel, data HlsVideoEncodingParams) error {
	body, _ := json.Marshal(data)

	err := ch.Publish("", HlsVideoEncodingQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}

func MediaEncodingCallbackProducer(ch Channel, data MediaEncodingCallbackParams) error {
	body, _ := json.Marshal(data)

	err := ch.Publish("", MediaEncodingCallbackQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
