package worker

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type VideoTranscodingParams struct {
	MediaUUID      uuid.UUID
	SourceFilePath string
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
