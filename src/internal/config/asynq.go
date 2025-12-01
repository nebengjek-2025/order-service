package config

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func NewAsynqClient(v *viper.Viper) *asynq.Client {
	host := v.GetString("redis.host")
	if host == "" {
		host = "127.0.0.1"
	}

	port := v.GetString("redis.port")
	addr := fmt.Sprintf("%s:%v", host, port)

	redisOpt := asynq.RedisClientOpt{
		Addr:     addr,
		Username: "default",
		Password: v.GetString("redis.password"),
		DB:       v.GetInt("redis.db"),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
		DialTimeout:  5 * time.Second,
		PoolSize:     10,
	}

	return asynq.NewClient(redisOpt)
}
