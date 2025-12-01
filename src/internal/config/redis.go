package config

import (
	redisModule "order-service/src/pkg/redis"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func LoadRedisConfig(viper *viper.Viper) {

	CfgRedis := &redisModule.CfgRedis{
		UseCluster:           viper.GetString("redis.use_cluster") == "true",
		EnableTLS:            true,
		RedisHost:            viper.GetString("redis.host"),
		RedisPort:            viper.GetString("redis.port"),
		RedisPassword:        viper.GetString("redis.password"),
		RedisDB:              viper.GetInt("redis.db"),
		RedisClusterNode:     viper.GetString("redis.cluster.node"),
		RedisClusterPassword: viper.GetString("redis.cluster.password"),
	}
	redisModule.LoadConfig(CfgRedis)
	redisModule.InitConnection()
}

func NewRedis() redis.UniversalClient {
	return redisModule.GetClient()
}
