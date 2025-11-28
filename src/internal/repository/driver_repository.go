package repository

import (
	"context"
	"fmt"
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

func (r *DriverRepository) GetDetailDriver(ctx context.Context, id string) (*entity.DriverInfo, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	var drivers entity.DriverInfo
	query := `
		SELECT 
			i.driver_id,
			u.full_name,
			i.city,
			i.jenis_kendaraan,
			i.nopol
		FROM users u
		JOIN info_driver i 
			ON u.user_id = i.driver_id
		WHERE u.user_id = ?;
		`

	err = db.GetContext(ctx, &drivers, query, id)
	fmt.Println(err, "<<<<BABI JOKOWI")
	if err != nil {
		return nil, err
	}

	return &drivers, nil
}

func (r *DriverRepository) SetOnTrip(ctx context.Context, driverID string) error {
	db, err := r.DB.GetDB()
	if err != nil {
		return fmt.Errorf("failed get db: %w", err)
	}

	query := `
		UPDATE driver_availability
		SET 
			is_available = 0,
			status = 'on_trip',
			last_seen_at = NOW()
		WHERE driver_id = ?
	`

	res, err := db.ExecContext(ctx, query, driverID)
	if err != nil {
		return fmt.Errorf("failed update driver_availability: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed get rows affected: %w", err)
	}

	if rows == 0 {
		insertQ := `
			INSERT INTO driver_availability (
				driver_id, is_available, status, last_seen_at
			) VALUES (?, 0, 'on_trip', NOW())
		`

		if _, err := db.ExecContext(ctx, insertQ, driverID); err != nil {
			return fmt.Errorf("failed insert driver_availability: %w", err)
		}
	}

	return nil
}

func (r *DriverRepository) SetOnline(ctx context.Context, driverID string) error {
	db, err := r.DB.GetDB()
	if err != nil {
		return err
	}

	q := `
		UPDATE driver_availability
		SET is_available = 1,
		    status = 'online',
		    last_seen_at = NOW()
		WHERE driver_id = ?
	`

	_, err = db.ExecContext(ctx, q, driverID)
	return err
}
