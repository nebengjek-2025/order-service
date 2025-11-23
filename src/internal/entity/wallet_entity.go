package entity

import "time"

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
