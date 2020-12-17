package worker

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func videoTranscodingConsumer(w *Worker, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var body VideoTranscodingParams

		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			w.logger.Error("Unmarshal error", zap.Error(err))
			return
		}

		w.logger.Info("Received a message", zap.String("SourceFilePath", body.SourceFilePath))

		// Fetch Media A
		// Create a new MediaFile for this Media
		//_, _ = w.storage.Open(context.Background(), "uuid/source")

		err = d.Ack(false)
		if err != nil {
			w.logger.Error("Error trying to send ack", zap.Error(err))
		}
	}
}
