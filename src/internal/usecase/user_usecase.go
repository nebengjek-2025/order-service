package usecase

import (
	"context"
	"fmt"
	"order-service/src/internal/model"
	"order-service/src/internal/model/converter"
	httpError "order-service/src/pkg/http-error"
	"order-service/src/pkg/utils"

	// "order-service/src/internal/gateway/messaging"

	"order-service/src/internal/repository"
	"order-service/src/pkg/log"

	"github.com/go-playground/validator/v10"
)

type UserUseCase struct {
	Log            log.Log
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(logger log.Log, validate *validator.Validate, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
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
