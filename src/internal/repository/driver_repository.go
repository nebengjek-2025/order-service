package repository

import (
	"context"
	"order-service/src/internal/entity"
	"order-service/src/pkg/databases/mysql"
)

type DriverRepository struct {
	DB mysql.DBInterface
}

func NewDriverRepository(db mysql.DBInterface) *DriverRepository {
	return &DriverRepository{
		DB: db,
	}
}

func (r *DriverRepository) FindDriver(ctx context.Context, id string) ([]entity.AvailableDriver, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	var drivers []entity.AvailableDriver
	query := `
		SELECT 
			da.driver_id,
			da.status,
			da.last_seen_at,
			i.city,
			i.jenis_kendaraan
		FROM driver_availability da
		JOIN info_driver i 
			ON da.driver_id = i.driver_id
		WHERE da.is_available = 1
		AND da.status = 'online'
		AND da.last_seen_at >= NOW() - INTERVAL 2 MINUTE
		AND da.driver_id = ?;
		`

	err = db.SelectContext(ctx, &drivers, query, id)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}
