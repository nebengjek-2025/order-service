package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient redis.UniversalClient

func InitConnection() {
	if !AppConfigData.UseCluster {
		var tlsConf *tls.Config
		if RedisConfigData.EnableTLS {
			tlsConf = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%v", RedisConfigData.Host, RedisConfigData.Port),
			Password:     RedisConfigData.Password,
			DB:           RedisConfigData.DB,
			TLSConfig:    tlsConf,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MaxRetries:   2,
		})

		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
			fmt.Println("REDIS ERROR:", err.Error())
			panic("cannot connect to Redis")
		}
	} else {
		var tlsConf *tls.Config
		if RedisClusterConfigData.EnableTLS {
			tlsConf = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}

		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        RedisClusterConfigData.Hosts,
			Username:     RedisClusterConfigData.Username,
			Password:     RedisClusterConfigData.Password,
			TLSConfig:    tlsConf,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
			fmt.Println("REDIS CLUSTER ERROR:", err.Error())
			panic("Cannot connect to Redis Cluster")
		}
	}
}

func GetClient() redis.UniversalClient {
	return redisClient
}
