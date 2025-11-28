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

func (p *Producer[T]) SendTo(topic string, event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		p.Log.Error("gateway/messaging/producer", "failed to marshal event", "SendTo", err.Error())
		return err
	}

	message := &k.Message{
		TopicPartition: k.TopicPartition{Topic: &topic, Partition: k.PartitionAny},
		Key:            []byte(event.GetId()),
		Value:          value,
	}

	err = p.Producer.Publish(message)
	if err != nil {
		p.Log.Error("gateway/messaging/producer", "error send message", "SendTo", err.Error())
		return err
	}

	p.Log.Info("gateway/messaging/producer", "event published successfully", topic, event.GetId())
	return nil
}

func (p *Producer[T]) Send(event T) error {
	return p.SendTo(p.Topic, event)
}
