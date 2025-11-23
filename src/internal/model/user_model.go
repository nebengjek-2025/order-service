package model

import "time"

type UserResponse struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	MobileNumber string     `json:"mobile_number,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
