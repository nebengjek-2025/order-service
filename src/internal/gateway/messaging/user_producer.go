package messaging

import (
	"order-service/src/internal/model"
	kafka "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"
)

type UserProducer struct {
	RequestRideProducer Producer[*model.UserEvent]
	DriverMatchProducer Producer[*model.DriverMatchEvent]
	Producer[*model.UserEvent]
}

func NewUserProducer(producer kafka.Producer, log log.Log) *UserProducer {
	return &UserProducer{
		RequestRideProducer: Producer[*model.UserEvent]{
			Producer: producer,
			Topic:    "request-ride",
			Log:      log,
		},
		DriverMatchProducer: Producer[*model.DriverMatchEvent]{
			Producer: producer,
			Topic:    "driver-match",
			Log:      log,
		},
	}
}

func (u *UserProducer) SendRequestRide(event *model.UserEvent) error {
	return u.RequestRideProducer.Send(event)
}

func (u *UserProducer) SendDriverMatch(event *model.DriverMatchEvent) error {
	return u.DriverMatchProducer.Send(event)
}
