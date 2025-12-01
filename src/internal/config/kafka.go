package config

import (
	kafkaPkgConfluent "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"

	"github.com/spf13/viper"
)

func NewKafkaConfig(viper *viper.Viper) kafkaPkgConfluent.KafkaConfig {
	configKafka := kafkaPkgConfluent.Cfg{
		KafkaUrl:      viper.GetString("kafka.bootstrap.servers"),
		KafkaUsername: viper.GetString("kafka.username"),
		KafkaPassword: viper.GetString("kafka.password"),
		KafkaCaCert:   viper.GetString("kafka.cacert"),
		AppName:       viper.GetString("kafka.app.name"),
	}
	return kafkaPkgConfluent.InitKafkaConfig(configKafka)

}

func NewKafkaProducer(config *viper.Viper, log log.Log) kafkaPkgConfluent.Producer {
	if !config.GetBool("kafka.producer.enabled") {
		log.Info("kafka-config", "Kafka producer is disabled in configuration", "kafka", "")
		return nil
	}
	kafkaProducer, err := kafkaPkgConfluent.NewProducer(kafkaPkgConfluent.GetConfig().GetKafkaConfig(), log)
	if err != nil {
		panic(err)
	}

	return kafkaProducer
}
