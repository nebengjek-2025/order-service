package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/src/internal/entity"
	"order-service/src/internal/gateway/messaging"
	"order-service/src/internal/model"
	"order-service/src/internal/model/converter"
	httpError "order-service/src/pkg/http-error"
	"order-service/src/pkg/log"
	"order-service/src/pkg/utils"
	"time"

	// "order-service/src/internal/gateway/messaging"

	"order-service/src/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type DriverUseCase struct {
	Log              log.Log
	Validate         *validator.Validate
	UserRepository   *repository.UserRepository
	WalletRepository *repository.WalletRepository
	OrderRepository  *repository.OrderRepository
	DriverRepository *repository.DriverRepository
	Config           *viper.Viper
	Redis            redis.UniversalClient
	DriverProducer   *messaging.DriverProducer
}

func NewDriverUseCase(
	logger log.Log,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	driverRepository *repository.DriverRepository,
	orderRepository *repository.OrderRepository,
	walletRepository *repository.WalletRepository,
	cfg *viper.Viper,
	redisClient redis.UniversalClient,
	driverProducer *messaging.DriverProducer,
) *DriverUseCase {
	return &DriverUseCase{
		Log:              logger,
		Validate:         validate,
		UserRepository:   userRepository,
		DriverRepository: driverRepository,
		OrderRepository:  orderRepository,
		WalletRepository: walletRepository,
		Config:           cfg,
		Redis:            redisClient,
		DriverProducer:   driverProducer,
	}
}

func (c *DriverUseCase) PickupPassanger(ctx context.Context, request *model.PickupPassanger) utils.Result {
	var result utils.Result

	if err := c.Validate.Struct(request); err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("validation error: %v", err.Error())
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", utils.ConvertString(err))
		return result
	}

	driverInfo, err := c.UserRepository.FindByID(ctx, request.DriverID)
	if err != nil {
		c.Log.Error("driver-usecase", fmt.Sprintf("Error get data driver :%v", err), "find-driver-info", utils.ConvertString(err))
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("driver with id %s not found", request.DriverID)
		result.Error = errObj
		return result
	}
	tripOrder, err := c.OrderRepository.FindOneOrder(ctx, entity.OrderFilter{PassengerID: &request.PassangerID, OrderID: &request.OrderID})
	if err != nil || tripOrder == nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Order not found"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", utils.ConvertString(err))
		return result
	}
	switch tripOrder.Status {
	case "COMPLETED", "CANCELLED":
		errObj := httpError.NewConflict()
		errObj.Message = "Order is no longer available (already completed or cancelled)"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", "")
		return result

	case "ON_GOING":
		errObj := httpError.NewConflict()
		errObj.Message = "Trip is already in progress"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", "")
		return result

	case "REQUESTED", "MATCHING":
		// Passenger belum confirm driver
		errObj := httpError.NewConflict()
		errObj.Message = "Passenger has not confirmed a driver yet for this order"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", "")
		return result

	case "ACCEPTED":
		// continue
	default:
		errObj := httpError.NewConflict()
		errObj.Message = fmt.Sprintf("Order in invalid state for pickup: %s", tripOrder.Status)
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", tripOrder.Status)
		return result
	}

	if tripOrder.DriverID == nil || *tripOrder.DriverID == "" {
		errObj := httpError.NewConflict()
		errObj.Message = "No driver assigned to this order yet"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", "")
		return result
	}

	if *tripOrder.DriverID != driverInfo.UserID {
		errObj := httpError.NewConflict()
		errObj.Message = "You are not the assigned driver for this order"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", "")
		return result
	}

	ok, err := c.OrderRepository.UpdateStatusOrderForDriver(ctx, request.OrderID, driverInfo.UserID, "ACCEPTED", "ON_GOING")
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "Failed to update order status to ON_GOING"
		result.Error = errObj
		c.Log.Error("driver-usecase", fmt.Sprintf("Error update status order: %v", err), "PickupPassanger", "")
		return result
	}

	if !ok {
		errObj := httpError.NewConflict()
		errObj.Message = "Order could not be updated to ON_GOING. It may have been changed or cancelled."
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", "concurrent-update")
		return result
	}
	tripOrder.Status = "ON_GOING"
	if err := c.DriverRepository.SetOnTrip(ctx, driverInfo.UserID); err != nil {
		c.Log.Error("driver-usecase", fmt.Sprintf("Failed update driver availability: %v", err), "PickupPassanger", "")
		// to do send event to handle failed update driver availableity
	}

	marshaledData, _ := json.Marshal(tripOrder)
	c.Log.Info("driver-usecase", "marshaled trip order", "PickupPassanger", utils.ConvertString(marshaledData))
	key := fmt.Sprintf("DRIVER:PICKING-PASSANGER:%s", driverInfo.UserID)
	redisErr := c.Redis.Set(ctx, key, marshaledData, 2*time.Hour).Err()
	if redisErr != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Internal server error insert to redis: %v", redisErr.Error())
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", utils.ConvertString(redisErr.Error()))
		return result
	}

	event := converter.OrderToEvent(request)
	c.Log.Info("driver-usecase", "Publishing user created event", "FindDriver", utils.ConvertString(event))
	if err = c.DriverProducer.Send(event); err != nil {
		c.Log.Error("driver-usecase", fmt.Sprintf("Failed publish driver created event : %+v", err), "order created", "")
	}

	result.Data = tripOrder

	return result
}
