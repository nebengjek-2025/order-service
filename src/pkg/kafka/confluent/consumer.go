package kafka

import (
	"fmt"
	"order-service/src/pkg/log"
	"strings"
	"sync"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type consumer struct {
	sync.Mutex
	handler  ConsumerHandler
	consumer *kafka.Consumer
	logger   log.Log
}

// NewConsumer is a constructor of kafka consumer
func NewConsumer(config *kafka.ConfigMap, log log.Log) (Consumer, error) {
	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	return &consumer{
		logger:   log,
		consumer: c,
	}, nil
}

func (c *consumer) SetHandler(handler ConsumerHandler) {
	c.handler = handler
}

func (c *consumer) Subscribe(topics ...string) {
	if c.handler == nil {
		joinTopic := strings.Join(topics, ", ")
		msg := fmt.Sprintf("Kafka Consumer Error: Topics: [%s] There is no consumer handler to handle message from incoming event", joinTopic)
		c.logger.Error("", msg, "", "")
		return
	}
	var wg sync.WaitGroup
	var mtx sync.Mutex

	c.consumer.SubscribeTopics(topics, nil)
	go func() {
		for {
			wg.Add(1)

			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				msg := fmt.Sprintf("Kafka Consumer Error: %v (%v)\n", err, msg)
				c.logger.Error("", msg, "", "")
				continue
			}
			mtx.Lock()
			c.handler.HandleMessage(msg)
			mtx.Unlock()
			wg.Done()
			c.consumer.CommitMessage(msg)
		}
	}()
	wg.Wait()
}
