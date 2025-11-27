package repository

import (
	"context"
	"order-service/src/internal/entity"
	"order-service/src/pkg/databases/mysql"
)

type OrderRepository struct {
	DB mysql.DBInterface
}

func NewOrderRepository(db mysql.DBInterface) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (r *OrderRepository) FindOrder(ctx context.Context, id string) (*entity.Order, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	var order entity.Order

	query := `
		SELECT 
			o.id,
			o.order_id,
			o.passenger_id,
			o.driver_id,
			o.created_at,
			o.updated_at
		FROM orders o
		WHERE o.order_id = ?
	`

	err = db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) OrderDetail(ctx context.Context, id string) (*entity.OrderDetail, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	var order entity.OrderDetail

	query := `
		SELECT 
			0.order_id,
			o.passenger_id,
			o.driver_id,
			o.origin_lat,
			o.origin_lng,
			o.destination_lat,
			o.destination_lng,
			o.origin_address,
			o.destination_address,
			o.route,
			o.min_price,
			o.max_price,
			o.best_route_km,
			o.best_route_price,
			o.best_route_duration,
			o.status,
			o.payment_method,
			o.payment_status,
			o.created_at,
			o.updated_at,

			pt.id AS payment_id,
			pt.amount AS payment_amount,
			pt.payment_status AS payment_status_detail,
			pt.provider_name AS payment_provider,
			pt.provider_reference_id AS payment_ref_id,
			pt.paid_at AS payment_paid_at,

			pr.id AS redemption_id,
			pr.discount_applied,
			pc.promo_code,
			pc.name AS promo_name,
			pc.discount_type,
			pc.discount_value,
			pc.max_discount
		FROM orders o
		LEFT JOIN payment_transactions pt ON pt.ride_order_id = o.order_id
		LEFT JOIN promo_redemptions pr ON pr.ride_order_id = o.order_id
		LEFT JOIN promo_campaigns pc ON pc.id = pr.promo_campaign_id
		WHERE o.order_id = ?
	`

	err = db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
