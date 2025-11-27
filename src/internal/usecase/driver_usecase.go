package usecase

import (
	"context"
	"order-service/src/internal/gateway/messaging"
	"order-service/src/internal/model"
	"order-service/src/pkg/log"
	"order-service/src/pkg/utils"

	// "order-service/src/internal/gateway/messaging"

	"order-service/src/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"googlemaps.github.io/maps"
)

type DriverUseCase struct {
	Log              log.Log
	Validate         *validator.Validate
	UserRepository   *repository.UserRepository
	WalletRepository *repository.WalletRepository
	OrderRepository  *repository.OrderRepository
	Config           *viper.Viper
	Redis            redis.UniversalClient
	UserProducer     *messaging.UserProducer
	Geoservice       *maps.Client
}

func NewDriverUseCase(
	logger log.Log,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	walletRepository *repository.WalletRepository,
	cfg *viper.Viper,
	redisClient redis.UniversalClient,
	userProducer *messaging.UserProducer,
	geo *maps.Client,
) *DriverUseCase {
	return &DriverUseCase{
		Log:              logger,
		Validate:         validate,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
		Config:           cfg,
		Redis:            redisClient,
		UserProducer:     userProducer,
		Geoservice:       geo,
	}
}

func (c *DriverUseCase) PickupPassanger(ctx context.Context, request *model.PickupPassanger) utils.Result {
	var result utils.Result

	return result
}
