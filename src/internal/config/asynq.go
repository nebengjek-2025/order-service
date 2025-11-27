package config

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func NewAsynqClient(v *viper.Viper) *asynq.Client {
	host := v.GetString("redis.host")
	if host == "" {
		host = "127.0.0.1"
	}

	port := v.GetInt("redis.port")
	if port == 0 {
		port = 6379
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	redisOpt := asynq.RedisClientOpt{
		Addr: addr,
	}

	return asynq.NewClient(redisOpt)
}
