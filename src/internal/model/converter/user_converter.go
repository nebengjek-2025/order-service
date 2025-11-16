package converter

import (
	"order-service/src/internal/entity"
	"order-service/src/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:           user.UserID,
		Name:         user.FullName,
		MobileNumber: user.MobileNumber,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
