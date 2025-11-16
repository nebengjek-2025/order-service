package main

import (
	"fmt"
	"order-service/src/internal/config"
	"order-service/src/internal/delivery/http/middleware"
	"order-service/src/pkg/log"
)

func main() {

	viperConfig := config.NewViper()
	viperConfig.SetDefault("log.level", "DEBUG")
	viperConfig.SetDefault("app.name", "ORDER_SERVICE")
	viperConfig.SetDefault("web.port", 8080)
	log.InitLogger(viperConfig)
	config.NewKafkaConfig(viperConfig)
	logger := log.GetLogger()
	config.LoadRedisConfig(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	redisClient := config.NewRedis()
	// producer := config.NewKafkaProducer(viperConfig, logger)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	app.Use(middleware.NewLogger())
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      logger,
		Validate: validate,
		Config:   viperConfig,
		// Producer:    producer,
		Redis: redisClient,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.GetLogger().Error("main", fmt.Sprintf("Failed to start server: %v", err), "main", "")
	}
}
