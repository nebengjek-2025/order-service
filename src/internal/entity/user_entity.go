package entity

import "time"

type User struct {
	UserID       string     `json:"user_id" db:"user_id"`
	FullName     string     `json:"full_name" db:"full_name"`
	Email        string     `json:"email" db:"email"`
	IsMitra      bool       `json:"isMitra" db:"isMitra"`
	MobileNumber string     `json:"mobile_number" db:"mobile_number"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}
