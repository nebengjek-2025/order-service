package kafka

import (
	"encoding/base64"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Producer interface {
	Publish(message *k.Message) error
	PublishChannel(topic string, message []byte)
}

type Consumer interface {
	SetHandler(handler ConsumerHandler)
	Subscribe(topics ...string)
}

type ConsumerHandler interface {
	HandleMessage(message *k.Message)
}

type KafkaConfig struct {
	Username      string
	Password      string
	Address       string
	SaslMechanism string
	AppName       string
	KafkaCaCert   string
}

type Cfg struct {
	KafkaUrl      string
	KafkaUsername string
	KafkaPassword string
	KafkaCaCert   string
	AppName       string
}

var kafkaConfig KafkaConfig

func InitKafkaConfig(cfg Cfg) KafkaConfig {

	kafkaConfig = KafkaConfig{
		Address:       cfg.KafkaUrl,
		Username:      cfg.KafkaUsername,
		Password:      cfg.KafkaPassword,
		AppName:       cfg.AppName,
		KafkaCaCert:   cfg.KafkaCaCert,
		SaslMechanism: "PLAIN",
	}
	return kafkaConfig
}

func GetConfig() KafkaConfig {
	return kafkaConfig
}

func decodeKey(secret string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func (kc KafkaConfig) GetKafkaConfig() *k.ConfigMap {
	kafkaCfg := k.ConfigMap{}

	if kc.Username != "" {
		ca, _ := decodeKey(kc.KafkaCaCert)
		kafkaCfg["sasl.mechanism"] = kc.SaslMechanism
		kafkaCfg["sasl.username"] = kc.Username
		kafkaCfg["sasl.password"] = kc.Password
		kafkaCfg["ssl.ca.pem"] = ca
		kafkaCfg["security.protocol"] = "sasl_ssl"
	}
	kafkaCfg.SetKey("bootstrap.servers", kc.Address)
	kafkaCfg.SetKey("group.id", kc.AppName)
	kafkaCfg.SetKey("retry.backoff.ms", 500)
	kafkaCfg.SetKey("socket.max.fails", 10)
	kafkaCfg.SetKey("reconnect.backoff.ms", 200)
	kafkaCfg.SetKey("reconnect.backoff.max.ms", 5000)
	kafkaCfg.SetKey("request.timeout.ms", 5000)
	kafkaCfg.SetKey("partition.assignment.strategy", "roundrobin")
	kafkaCfg.SetKey("auto.offset.reset", "earliest")

	return &kafkaCfg
}
