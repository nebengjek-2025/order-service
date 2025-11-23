package config

import (
	"order-service/src/internal/delivery/http"
	"order-service/src/internal/delivery/http/middleware"
	"order-service/src/internal/delivery/http/route"
	"order-service/src/internal/gateway/messaging"

	// "order-service/src/internal/gateway/messaging"
	"order-service/src/internal/repository"
	"order-service/src/internal/usecase"
	"order-service/src/pkg/databases/mysql"
	kafkaPkgConfluent "order-service/src/pkg/kafka/confluent"
	"order-service/src/pkg/log"

	"github.com/redis/go-redis/v9"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       mysql.DBInterface
	App      *fiber.App
	Log      log.Log
	Validate *validator.Validate
	Config   *viper.Viper
	Producer kafkaPkgConfluent.Producer
	Redis    redis.UniversalClient
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.DB)
	walletRepository := repository.NewWalletRepository(config.DB)
	orderRepository := repository.NewOrderRepository(config.DB)

	userProducer := messaging.NewUserProducer(config.Producer, config.Log)
	// setup use cases
	userUseCase := usecase.NewUserUseCase(
		config.Log,
		config.Validate,
		userRepository,
		walletRepository,
		orderRepository,
		config.Config,
		config.Redis,
		userProducer)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.VerifyBearer(config.Config)

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}
