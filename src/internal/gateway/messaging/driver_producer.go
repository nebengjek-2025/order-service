package messaging

import (
	"order-service/src/internal/model"
	kafka "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"
)

type DriverProducer struct {
	Producer[*model.OrderEvent]
}

func NewDriverProducer(producer kafka.Producer, log log.Log) *DriverProducer {
	return &DriverProducer{
		Producer: Producer[*model.OrderEvent]{
			Producer: producer,
			Topic:    "trip-created",
			Log:      log,
		},
	}
}
