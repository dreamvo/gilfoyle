package worker

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type EncodingEntrypointParams struct {
	MediaUUID uuid.UUID `json:"media_uuid"`
}

type HlsVideoEncodingParams struct {
	MediaFileUUID      uuid.UUID `json:"media_file_uuid"`
	KeyframeInterval   int       `json:"keyframe_interval"`
	HlsSegmentDuration int       `json:"hls_segment_duration"`
	HlsPlaylistType    string    `json:"hls_playlist_type"`
}

type EncodingFinalizerParams struct {
	MediaUUID uuid.UUID `json:"media_uuid"`
}

func EncodingEntrypointProducer(ch Channel, data EncodingEntrypointParams) error {
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

func EncodingFinalizerProducer(ch Channel, data EncodingFinalizerParams) error {
	body, _ := json.Marshal(data)

	err := ch.Publish("", EncodingFinalizerQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
