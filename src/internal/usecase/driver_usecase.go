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
	"strconv"
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

func (c *DriverUseCase) CompletedTrip(ctx context.Context, request *model.RequestCompleteTrip) utils.Result {
	var result utils.Result

	if err := c.Validate.Struct(request); err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("validation error: %v", err.Error())
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", utils.ConvertString(err))
		return result
	}

	tripOrder, err := c.OrderRepository.FindOneOrder(ctx, entity.OrderFilter{DriverID: &request.DriverID, OrderID: &request.OrderID})
	if err != nil || tripOrder == nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Order not found"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "PickupPassanger", utils.ConvertString(err))
		return result
	}
	if tripOrder.Status != "ON_GOING" {
		errObj := httpError.NewConflict()
		errObj.Message = fmt.Sprintf("Cannot complete trip in status %s", tripOrder.Status)
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", tripOrder.Status)
		return result
	}
	var tracker model.TripTracker
	key := fmt.Sprintf("trip:%s", request.OrderID)
	driverTracker, errRedis := c.Redis.Get(ctx, key).Result()
	if errRedis != nil || driverTracker == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("Error get data from redis: %v", errRedis)
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", utils.ConvertString(errRedis))
		return result
	}
	if err := json.Unmarshal([]byte(driverTracker), &tracker); err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error unmarshal tripdata: %v", err)
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", utils.ConvertString(err))
		return result
	}

	if tracker.Data.DriverID != "" && tracker.Data.DriverID != request.DriverID {
		errObj := httpError.NewConflict()
		errObj.Message = "Trip data does not belong to this driver"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", utils.ConvertString(tracker))
		return result
	}

	realDistance, err := strconv.ParseFloat(tracker.Data.Distance, 64)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Invalid distance value: %v", err)
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", tracker.Data.Distance)
		return result
	}

	if tripOrder.DriverID != nil && *tripOrder.DriverID != "" {
		keyStatusDriver := fmt.Sprintf("DRIVER:PICKING-PASSANGER:%s", *tripOrder.DriverID)
		if err := c.Redis.Del(ctx, keyStatusDriver).Err(); err != nil {
			c.Log.Error("driver-usecase", fmt.Sprintf("failed delete redis key %s: %v", keyStatusDriver, err), "CompletedTrip", "")
		}
	}
	duration := time.Since(tripOrder.CreatedAt)
	durationMinutes := int(duration.Minutes())
	durationFormatted := utils.FormatDuration(durationMinutes)
	ok, err := c.OrderRepository.CompleteTrip(ctx, request.OrderID, request.DriverID, realDistance, durationFormatted)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "Failed to complete trip"
		result.Error = errObj
		c.Log.Error("driver-usecase", fmt.Sprintf("Error updating order to COMPLETED: %v", err), "CompletedTrip", "")
		return result
	}
	if !ok {
		errObj := httpError.NewConflict()
		errObj.Message = "Order could not be completed, it may have been updated or cancelled"
		result.Error = errObj
		c.Log.Error("driver-usecase", errObj.Message, "CompletedTrip", "concurrent-update")
		return result
	}

	if err := c.DriverRepository.SetOnline(ctx, request.DriverID); err != nil {
		c.Log.Error("driver-usecase",
			fmt.Sprintf("Failed update driver availability to online: %v", err),
			"CompletedTrip",
			"")
	}

	_ = c.Redis.Del(ctx, key).Err()
	_ = c.Redis.Del(ctx, fmt.Sprintf("order:%s:distance", request.OrderID)).Err()
	_ = c.Redis.Del(ctx, fmt.Sprintf("order:%s:driver:%s", request.OrderID, request.DriverID)).Err()

	result.Data = map[string]interface{}{
		"order_id":        request.OrderID,
		"driver_id":       request.DriverID,
		"status":          "COMPLETED",
		"distance_actual": realDistance,
		"message":         "Trip completed successfully",
	}

	return result
}
