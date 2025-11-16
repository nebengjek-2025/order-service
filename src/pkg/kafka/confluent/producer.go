package kafka

import (
	"fmt"
	"order-service/src/pkg/log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer struct
type producer struct {
	producer *kafka.Producer
	logger   log.Log
}

// NewProducer constructor
func NewProducer(config *kafka.ConfigMap, log log.Log) (Producer, error) {
	p, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	prod := &producer{
		producer: p,
		logger:   log,
	}

	go prod.errReporter()

	return prod, nil
}

func (p *producer) errReporter() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				msg := fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition)
				p.logger.Error("", msg, "", "")
			}
		}
	}
}

func (p *producer) Publish(topic string, message []byte) {
	msgCh := p.producer.ProduceChannel()
	msgCh <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}
}
