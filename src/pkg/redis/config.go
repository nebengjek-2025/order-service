package redis

import (
	"fmt"
	"order-service/src/pkg/utils"
	"strings"
)

type CfgRedis struct {
	UseRedis             bool
	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RedisDB              int
	RedisClusterNode     string
	RedisClusterPassword string
}

type AppConfig struct {
	UseRedis bool
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RedisClusterConfig struct {
	Hosts    []string
	Password string
}

var (
	AppConfigData          AppConfig
	RedisConfigData        RedisConfig
	RedisClusterConfigData RedisClusterConfig
)

func LoadConfig(config *CfgRedis) {

	AppConfigData = AppConfig{
		UseRedis: config.UseRedis,
	}

	redisDb := config.RedisDB
	redisHost := config.RedisHost
	redisPort := config.RedisPort
	redisPass := config.RedisPassword

	RedisConfigData = RedisConfig{
		Host:     fmt.Sprintf("%v", redisHost),
		Port:     fmt.Sprintf("%v", redisPort),
		Password: fmt.Sprintf("%v", redisPass),
		DB:       utils.ConvertInt(redisDb),
	}

	clusterHost := strings.Split(config.RedisClusterNode, ";")
	RedisClusterConfigData = RedisClusterConfig{
		Hosts:    clusterHost,
		Password: config.RedisClusterPassword,
	}
}
