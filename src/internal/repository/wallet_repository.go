package repository

import (
	"context"
	"database/sql"
	"order-service/src/internal/entity"
	"order-service/src/pkg/databases/mysql"
)

type WalletRepository struct {
	DB mysql.DBInterface
}

func NewWalletRepository(db mysql.DBInterface) *WalletRepository {
	return &WalletRepository{
		DB: db,
	}
}

func (r *WalletRepository) GetWalletByUserID(ctx context.Context, userID string) (*entity.Wallet, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	var w entity.Wallet
	query := `
		SELECT id, user_id, balance, last_updated, created_at, updated_at
		FROM wallets
		WHERE user_id = ?
		LIMIT 1
	`
	if err := db.GetContext(ctx, &w, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *WalletRepository) TopUp(ctx context.Context, userID string, amount float64) error {
	db, err := r.DB.GetDB()
	if err != nil {
		return err
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock wallet
	var wallet entity.Wallet
	err = tx.GetContext(ctx, &wallet, `
		SELECT id, balance FROM wallets
		WHERE user_id = ?
		FOR UPDATE
	`, userID)
	if err != nil {
		return err
	}

	// Update balance
	newBalance := wallet.Balance + amount
	_, err = tx.ExecContext(ctx, `
		UPDATE wallets SET balance = ?, last_updated = NOW(6)
		WHERE id = ?
	`, newBalance, wallet.ID)
	if err != nil {
		return err
	}

	// Insert transaction log
	_, err = tx.ExecContext(ctx, `
		INSERT INTO wallet_transactions
			(wallet_id, transaction_id, amount, type, description)
		VALUES (?, UUID(), ?, 'credit', 'Top up')
	`, wallet.ID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}
