package worker

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type VideoTranscodingParams struct {
	MediaUUID uuid.UUID
}

func ProduceVideoTranscodingQueue(ch *amqp.Channel, data VideoTranscodingParams) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ch.Publish("", VideoTranscodingQueue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
