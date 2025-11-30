package messaging

import (
	"order-service/src/internal/model"
	kafka "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"
)

type DriverProducer struct {
	DriverPickupProducer Producer[*model.OrderEvent]
	DriverUpdateProducer Producer[*model.NotificationUser]
	Producer[*model.OrderEvent]
}

func NewDriverProducer(producer kafka.Producer, log log.Log) *DriverProducer {
	return &DriverProducer{
		DriverPickupProducer: Producer[*model.OrderEvent]{
			Producer: producer,
			Topic:    "trip-created",
			Log:      log,
		},
		DriverUpdateProducer: Producer[*model.NotificationUser]{
			Producer: producer,
			Topic:    "order-driver-request-pickup",
			Log:      log,
		},
	}
}

func (u *DriverProducer) SendRequestRide(event *model.OrderEvent) error {
	return u.DriverPickupProducer.Send(event)
}

func (u *DriverProducer) SendOrderCompleted(event *model.NotificationUser) error {
	return u.DriverUpdateProducer.Send(event)
}
