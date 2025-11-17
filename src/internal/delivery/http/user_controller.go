package http

import (
	"order-service/src/internal/delivery/http/middleware"
	"order-service/src/internal/model"
	"order-service/src/internal/usecase"
	"order-service/src/pkg/log"
	"order-service/src/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Log     log.Log
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger log.Log) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.UserID,
	}
	result := c.UseCase.GetUser(ctx.Context(), request)
	if result.Error != nil {
		return utils.ResponseError(result.Error, ctx)
	}

	return utils.Response(result.Data, "GetProfile", fiber.StatusOK, ctx)
}

func (c *UserController) PostLocation(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.LocationSuggestionRequest)
	request.UserID = auth.UserID
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Error("UserController.PostLocation", "Failed to parse request body", "error", err.Error())
		return utils.ResponseError(err, ctx)
	}
	result := c.UseCase.PostLocation(ctx.Context(), request)
	if result.Error != nil {
		return utils.ResponseError(result.Error, ctx)
	}

	return utils.Response(result.Data, "Location Suggestion", fiber.StatusOK, ctx)
}

func (c *UserController) FindDriver(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.FindDriverRequest)
	request.UserID = auth.UserID
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Error("UserController.FindDriver", "Failed to parse request body", "error", err.Error())
		return utils.ResponseError(err, ctx)
	}
	result := c.UseCase.FindDriver(ctx.Context(), request)
	if result.Error != nil {
		return utils.ResponseError(result.Error, ctx)
	}

	return utils.Response(result.Data, "Find Driver", fiber.StatusOK, ctx)
}
