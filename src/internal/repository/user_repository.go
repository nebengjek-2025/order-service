package repository

import (
	"context"
	"order-service/src/internal/entity"
	"order-service/src/pkg/databases/mysql"
)

type UserRepository struct {
	DB mysql.DBInterface
}

func NewUserRepository(db mysql.DBInterface) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT user_id, full_name,mobile_number, isMitra, created_at, updated_at
		FROM users
		WHERE user_id = ?`

	err = db.GetContext(ctx, &user, query, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
