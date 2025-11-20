package redis

import (
	"fmt"
	"order-service/src/pkg/utils"
	"strings"
)

type CfgRedis struct {
	UseCluster           bool
	EnableTLS            bool
	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RedisDB              int
	RedisClusterNode     string
	RedisClusterPassword string
}

type AppConfig struct {
	UseCluster bool
}

type RedisConfig struct {
	Host      string
	Port      string
	Password  string
	DB        int
	EnableTLS bool
}

type RedisClusterConfig struct {
	Hosts     []string
	Username  string
	Password  string
	EnableTLS bool
}

var (
	AppConfigData          AppConfig
	RedisConfigData        RedisConfig
	RedisClusterConfigData RedisClusterConfig
)

func LoadConfig(config *CfgRedis) {

	AppConfigData = AppConfig{
		UseCluster: config.UseCluster,
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
