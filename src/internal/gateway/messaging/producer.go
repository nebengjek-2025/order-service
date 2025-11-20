package messaging

import (
	"encoding/json"
	"order-service/src/internal/model"
	kafka "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Producer[T model.Event] struct {
	Producer kafka.Producer
	Topic    string
	Log      log.Log
}

func (p *Producer[T]) GetTopic() *string {
	return &p.Topic
}

func (p *Producer[T]) Send(event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		p.Log.Error("gateway/messaging/producer", "failed to marshal event", "Send", err.Error())
		return err
	}

	message := &k.Message{
		TopicPartition: k.TopicPartition{Topic: &p.Topic, Partition: k.PartitionAny},
		Key:            []byte(event.GetId()),
		Value:          value,
	}

	err = p.Producer.Publish(message)
	if err != nil {
		p.Log.Error("send-event", "error send message", "send", err.Error())
		return err
	}

	return nil
}
