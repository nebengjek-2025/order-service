package entity

import "time"

type Wallet struct {
	ID          string    `db:"id"        json:"id"`
	UserID      string    `db:"user_id"   json:"user_id"`
	Balance     float64   `db:"balance"   json:"balance"`
	LastUpdated time.Time `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updated_at"`
}

type WalletTransaction struct {
	ID            uint64    `db:"id"             json:"id"`
	WalletID      string    `db:"wallet_id"      json:"wallet_id"`
	TransactionID string    `db:"transaction_id" json:"transaction_id"`
	Amount        float64   `db:"amount"         json:"amount"`
	Type          string    `db:"type"           json:"type"`
	Description   string    `db:"description"    json:"description"`
	Timestamp     time.Time `db:"timestamp"      json:"timestamp"`
	CreatedAt     time.Time `db:"created_at"     json:"created_at"`
}
