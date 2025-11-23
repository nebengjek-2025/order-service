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
	producer := config.NewKafkaProducer(viperConfig, logger)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	app.Use(middleware.NewLogger())
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      logger,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
		Redis:    redisClient,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.GetLogger().Error("main", fmt.Sprintf("Failed to start server: %v", err), "main", "")
	}
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

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
