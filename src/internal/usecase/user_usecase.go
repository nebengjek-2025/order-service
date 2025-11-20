package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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
	"googlemaps.github.io/maps"
)

type UserUseCase struct {
	Log              log.Log
	Validate         *validator.Validate
	UserRepository   *repository.UserRepository
	WalletRepository *repository.WalletRepository
	Config           *viper.Viper
	Redis            redis.UniversalClient
	UserProducer     *messaging.UserProducer
}

func NewUserUseCase(
	logger log.Log,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	walletRepository *repository.WalletRepository,
	cfg *viper.Viper,
	redisClient redis.UniversalClient,
	userProducer *messaging.UserProducer,
) *UserUseCase {
	return &UserUseCase{
		Log:              logger,
		Validate:         validate,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
		Config:           cfg,
		Redis:            redisClient,
		UserProducer:     userProducer,
	}
}

func (c *UserUseCase) GetUser(ctx context.Context, request *model.GetUserRequest) utils.Result {
	var result utils.Result

	if err := c.Validate.Struct(request); err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("validation error: %v", err.Error())
		result.Error = errObj
		c.Log.Error("GetUser-validation", err.Error(), "request", utils.ConvertString(request))
		return result
	}
	user, err := c.UserRepository.FindByID(ctx, request.ID)
	fmt.Println(err)
	if err != nil {
		c.Log.Error("GetUser-FindByID", err.Error(), "request", utils.ConvertString(request))
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("user with id %s not found", request.ID)
		result.Error = errObj
		return result
	}
	c.Log.Info("GetUser", "user found", "userID", request.ID)
	result.Data = converter.UserToResponse(user)
	return result
}

func (c *UserUseCase) PostLocation(ctx context.Context, request *model.LocationSuggestionRequest) utils.Result {
	var result utils.Result
	mapsClient, err := maps.NewClient(maps.WithAPIKey(c.Config.GetString("thirdparty.google.api_key")))
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("error creating Google Maps client: %v", err)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "init googleclient", utils.ConvertString(err))
		return result
	}

	routeSuggestion, err := c.getRouteSuggestions(ctx, mapsClient, request.CurrentLocation, request.Destination)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("error getRouteSuggestions: %v", err)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "PostLocation", utils.ConvertString(err))
		return result
	}

	key := fmt.Sprintf("USER:ROUTE:%s", request.UserID)
	routeSuggestion.Route.Origin = request.CurrentLocation
	routeSuggestion.Route.Destination = request.Destination
	routeSummaryJSON, err := json.Marshal(routeSuggestion)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error marshalling RouteSummary: %v", err)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "PostLocation", utils.ConvertString(err))
		return result
	}

	redisErr := c.Redis.Set(ctx, key, routeSummaryJSON, 60*time.Minute).Err()
	if redisErr != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error saving to redis: %v", redisErr)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "PostLocation", utils.ConvertString(redisErr))
		return result
	}
	result.Data = routeSuggestion

	return result
}

func (c *UserUseCase) FindDriver(ctx context.Context, request *model.FindDriverRequest) utils.Result {
	var result utils.Result
	key := fmt.Sprintf("USER:ROUTE:%s", request.UserID)
	var tripPlan model.RouteSummary
	redisData, errRedis := c.Redis.Get(ctx, key).Result()
	if errRedis != nil || redisData == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("Error get data from redis: %v", errRedis)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "FindDriver", utils.ConvertString(errRedis))
		return result
	}
	err := json.Unmarshal([]byte(redisData), &tripPlan)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error unmarshal tripdata: %v", err)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "FindDriver", utils.ConvertString(err))
		return result
	}
	// check payment method request
	switch request.PaymentMethod {
	case "WALLET":
		// proceed
		walletCheck, err := c.WalletRepository.WalletCheck(ctx, request.UserID)
		if err != nil {
			errObj := httpError.NewInternalServerError()
			errObj.Message = fmt.Sprintf("Wallet not found: %v, Please create wallet first", err)
			result.Error = errObj
			c.Log.Error("user-usecase", errObj.Message, "FindDriver", utils.ConvertString(err))
			return result
		}
		if walletCheck.Balance <= tripPlan.MaxPrice {
			errObj := httpError.NewBadRequest()
			errObj.Message = "insufficient balance, please topup"
			result.Error = errObj
			c.Log.Error("user-usecase", errObj.Message, "FindDriver", "")
			return result
		}
	case "qris", "CASH":
		if tripPlan.MaxPrice < 1000 {
			errObj := httpError.NewBadRequest()
			errObj.Message = "minimum payment amount is 1,000"
			result.Error = errObj
			return result
		}

		if tripPlan.MaxPrice > 10000000 {
			errObj := httpError.NewBadRequest()
			errObj.Message = "maximum payment amount exceeded (10,000,000)"
			result.Error = errObj
			return result
		}
	default:
		errObj := httpError.NewBadRequest()
		errObj.Message = "invalid payment method, only 'WALLET' or 'CASH' allowed"
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "FindDriver", request.PaymentMethod)
		return result
	}

	radius := 3.0
	drivers, err := c.Redis.GeoRadius(ctx, "drivers-locations", tripPlan.Route.Origin.Longitude, tripPlan.Route.Origin.Latitude, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithDist:  true,
		WithCoord: true,
		Sort:      "ASC",
	}).Result()

	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error searching drivers: %v", err)
		result.Error = errObj
		c.Log.Error("user-usecase", errObj.Message, "FindDriver", utils.ConvertString(err))
		return result
	}
	posibleDriver := "No driver available. Don't worry, please try again later."
	if len(drivers) > 0 { // c.Producer.Publish("request-ride", marshaledData)
		event := converter.UserToEvent(&model.RequestRide{
			UserId:       request.UserID,
			RouteSummary: tripPlan,
		})
		c.Log.Info("user-usecase", "Publishing user created event", "FindDriver", utils.ConvertString(event))
		if err = c.UserProducer.Send(event); err != nil {
			c.Log.Error("user-usecase", fmt.Sprintf("Failed publish user created event : %+v", err), "FindDriver", "")
			result.Error = httpError.NewInternalServerError()
			return result
		}
		posibleDriver = fmt.Sprintf("Please sit back, there are %d drivers available, we will let you know", len(drivers))
	}
	result.Data = model.FindDriverResponse{
		Message: posibleDriver,
		Driver:  drivers,
	}

	return result
}

func (c *UserUseCase) getRouteSuggestions(ctx context.Context, mapsClient *maps.Client, currentRequest model.LocationRequest, destinationRequest model.LocationRequest) (*model.RouteSummary, error) {
	origin := fmt.Sprintf("%f,%f", currentRequest.Latitude, currentRequest.Longitude)
	destination := fmt.Sprintf("%f,%f", destinationRequest.Latitude, destinationRequest.Longitude)
	departureTime := time.Now().Add(5 * time.Minute).Unix()

	req := &maps.DirectionsRequest{
		Origin:        origin,
		Destination:   destination,
		Mode:          maps.TravelModeDriving,
		Alternatives:  true,
		Optimize:      true,
		DepartureTime: fmt.Sprintf("%d", departureTime),
		TrafficModel:  maps.TrafficModelBestGuess,
	}

	routes, _, err := mapsClient.Directions(ctx, req)
	if err != nil {
		c.Log.Error("user-usecase", err.Error(), "getRouteSuggestions", fmt.Sprintf("Origin: %s, Destination: %s, err: %w", origin, destination, err.Error()))
		return nil, fmt.Errorf("error making directions request: %w", err)
	}

	if len(routes) == 0 {
		c.Log.Error("user-usecase", "no routes found", "getRouteSuggestions", fmt.Sprintf("Origin: %s, Destination: %s, result: %s", origin, destination, utils.ConvertString(routes)))
		return nil, fmt.Errorf("no routes found")
	}

	const pricePerKm = 3000.0
	var minPrice, maxPrice float64
	var bestRouteKm, bestRoutePrice, bestRouteDuration float64

	minPrice = math.MaxFloat64
	maxPrice = -math.MaxFloat64

	for _, route := range routes {
		totalDistance := 0.0
		totalDuration := 0.0

		for _, leg := range route.Legs {
			totalDistance += float64(leg.Distance.Meters)
			totalDuration += float64(leg.DurationInTraffic.Seconds())
		}

		distanceInKm := totalDistance / 1000.0
		price := distanceInKm * pricePerKm

		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}

		if bestRouteKm == 0 || price < bestRoutePrice {
			bestRouteKm = distanceInKm
			bestRoutePrice = price
			bestRouteDuration = totalDuration / 60
		}
	}

	return &model.RouteSummary{
		MinPrice:          minPrice,
		MaxPrice:          maxPrice,
		BestRouteKm:       bestRouteKm,
		BestRoutePrice:    bestRoutePrice,
		BestRouteDuration: utils.FormatDuration(int(math.Ceil(bestRouteDuration))),
		Duration:          int(math.Ceil(bestRouteDuration)),
	}, nil

}
