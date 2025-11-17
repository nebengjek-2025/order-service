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

type Wallet struct {
	ID             string           `bson:"_id,omitempty" json:"id"`
	UserID         string           `bson:"userId" json:"userId"`
	Balance        float64          `bson:"balance" json:"balance"`
	TransactionLog []TransactionLog `bson:"transactionLog" json:"transactionLog"`
	LastUpdated    time.Time        `bson:"lastUpdated" json:"lastUpdated"`
}

type TransactionLog struct {
	TransactionID string    `bson:"transactionId" json:"transactionId"`
	Amount        float64   `bson:"amount" json:"amount"`
	Type          string    `bson:"type" json:"type"`
	Description   string    `bson:"description" json:"description"`
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
}
