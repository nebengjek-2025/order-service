package main

import (
	"context"
	"fmt"
	"order-service/src/internal/config"
	"order-service/src/internal/delivery/http/middleware"
	"order-service/src/pkg/log"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func main() {

	viperConfig := config.NewViper()
	viperConfig.SetDefault("log.level", "DEBUG")
	viperConfig.SetDefault("app.name", "ORDER_SERVICE")
	viperConfig.SetDefault("web.port", 8080)
	log.InitLogger(viperConfig)
	config.NewKafkaConfig(viperConfig)
	logger := log.GetLogger()
	asynqClient := config.NewAsynqClient(viperConfig)
	config.LoadRedisConfig(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	redisClient := config.NewRedis()
	producer := config.NewKafkaProducer(viperConfig, logger)
	validate := config.NewValidator(viperConfig)
	geoservice, errG := config.NewGeoService(viperConfig)
	if errG != nil {
		logger.Error("main", fmt.Sprintf("Failed to initialize GeoService: %v", errG), "main", "")
		return
	}
	app := config.NewFiber(viperConfig)
	app.Use(middleware.NewLogger())
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%v", viperConfig.GetString("redis.host"), viperConfig.GetString("redis.port")),
		DB:   viper.GetInt("redis.db"),
	}
	if password := viper.GetString("redis.password"); password != "" {
		redisOpt.Password = password
	}

	asynqServer := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         logger,
		Validate:    validate,
		Config:      viperConfig,
		Producer:    producer,
		Redis:       redisClient,
		Geoservice:  geoservice,
		AsynqClient: asynqClient,
		Async:       mux,
	})
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	webPort := viperConfig.GetInt("web.port")
	go func() {
		logger.Info("main", fmt.Sprintf("Fiber listening on :%d", webPort), "fiber", "")
		if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
			logger.Error("main", fmt.Sprintf("Failed to start server: %v", err), "fiber", "")
		}
	}()

	go func() {
		logger.Info("main", "Asynq worker started", "asynq", "")
		if err := asynqServer.Run(mux); err != nil {
			logger.Error("main", fmt.Sprintf("Asynq server stopped with error: %v", err), "asynq", "")
		}
	}()

	go func() {
		<-quit
		logger.Info("main", "Server order-service is shutting down...", "gracefull", "")

		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.Shutdown(); err != nil {
			logger.Error("main", fmt.Sprintf("Error during shutdown: %v", err), "graceful", "")
		}
		close(done)
	}()

	<-done
	logger.Info("main", fmt.Sprintf("Server %s stopped", viperConfig.GetString("app.name")), "gracefull", "")
}
