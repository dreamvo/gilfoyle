package worker

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func videoTranscodingQueueConsumer(w *Worker, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var body VideoTranscodingParams

		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			w.Logger.Error("Unmarshal error", zap.Error(err))
		}

		w.Logger.Info("Received a message", zap.String("SourceFilePath", body.SourceFilePath))
	}
}
