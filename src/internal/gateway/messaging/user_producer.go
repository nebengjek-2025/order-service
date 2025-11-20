package messaging

import (
	"order-service/src/internal/model"
	kafka "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"
)

type UserProducer struct {
	Producer[*model.UserEvent]
}

func NewUserProducer(producer kafka.Producer, log log.Log) *UserProducer {
	return &UserProducer{
		Producer: Producer[*model.UserEvent]{
			Producer: producer,
			Topic:    "request-ride",
			Log:      log,
		},
	}
}
