package http

import (
	"order-service/src/internal/delivery/http/middleware"
	"order-service/src/internal/model"
	"order-service/src/internal/usecase"
	"order-service/src/pkg/log"
	"order-service/src/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type DriverController struct {
	Log     log.Log
	UseCase *usecase.DriverUseCase
}

func NewDriverController(useCase *usecase.DriverUseCase, logger log.Log) *DriverController {
	return &DriverController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *DriverController) PickupPassanger(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.PickupPassanger)
	request.DriverID = auth.UserID
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Error("UserController.PickupPassanger", "Failed to parse request body", "error", err.Error())
		return utils.ResponseError(err, ctx)
	}
	result := c.UseCase.PickupPassanger(ctx.Context(), request)
	if result.Error != nil {
		return utils.ResponseError(result.Error, ctx)
	}

	return utils.Response(result.Data, "Pickup Passanger", fiber.StatusOK, ctx)
}

func (c *DriverController) CompletedTrip(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.RequestCompleteTrip)
	request.DriverID = auth.UserID
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Error("UserController.CompletedTrip", "Failed to parse request body", "error", err.Error())
		return utils.ResponseError(err, ctx)
	}
	result := c.UseCase.CompletedTrip(ctx.Context(), request)
	if result.Error != nil {
		return utils.ResponseError(result.Error, ctx)
	}

	return utils.Response(result.Data, "Complete Trip", fiber.StatusOK, ctx)
}
