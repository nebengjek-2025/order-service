package repository

import (
	"context"
	"database/sql"
	"fmt"
	"order-service/src/internal/entity"
	"order-service/src/pkg/databases/mysql"
	"strings"
)

type OrderRepository struct {
	DB mysql.DBInterface
}

func NewOrderRepository(db mysql.DBInterface) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (r *OrderRepository) FindOrders(ctx context.Context, f entity.OrderFilter) ([]entity.Order, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return nil, err
	}

	baseQuery := `
		SELECT 
			o.id,
			o.order_id,
			o.passenger_id,
			o.driver_id,
			o.origin_lat,
			o.origin_lng,
			o.destination_lat,
			o.destination_lng,
			o.origin_address,
			o.destination_address,
			o.min_price,
			o.max_price,
			o.best_route_km,
			o.best_route_price,
			o.best_route_duration,
			o.status,
			o.payment_method,
			o.payment_status,
			o.estimated_fare,
			o.distance_km,
			o.distance_actual,
			o.duration_actual,
			o.created_at,
			o.updated_at
		FROM orders o
	`

	var (
		conds []string
		args  []interface{}
	)

	if f.OrderID != nil {
		conds = append(conds, "o.order_id = ?")
		args = append(args, *f.OrderID)
	}
	if f.PassengerID != nil {
		conds = append(conds, "o.passenger_id = ?")
		args = append(args, *f.PassengerID)
	}
	if f.DriverID != nil {
		conds = append(conds, "o.driver_id = ?")
		args = append(args, *f.DriverID)
	}
	if f.Status != nil {
		conds = append(conds, "o.status = ?")
		args = append(args, *f.Status)
	}

	if f.StatusNot != nil {
		conds = append(conds, "o.status != ?")
		args = append(args, *f.StatusNot)
	}

	if len(f.StatusIn) > 0 {
		placeholders := make([]string, 0, len(f.StatusIn))
		for range f.StatusIn {
			placeholders = append(placeholders, "?")
		}
		conds = append(conds, fmt.Sprintf("o.status IN (%s)", strings.Join(placeholders, ", ")))
		for _, s := range f.StatusIn {
			args = append(args, s)
		}
	}

	if f.PaymentStatus != nil {
		conds = append(conds, "o.payment_status = ?")
		args = append(args, *f.PaymentStatus)
	}

	query := baseQuery

	if len(conds) > 0 {
		query = query + " WHERE " + strings.Join(conds, " AND ")
	}

	query = query + " ORDER BY o.created_at DESC"

	var orders []entity.Order
	if err := db.SelectContext(ctx, &orders, query, args...); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) FindOneOrder(ctx context.Context, f entity.OrderFilter) (*entity.Order, error) {
	orders, err := r.FindOrders(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("order not found")
	}
	return &orders[0], nil
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

func (r *OrderRepository) InsertOrder(ctx context.Context, order *entity.CreateOrder) error {
	db, err := r.DB.GetDB()
	if err != nil {
		return fmt.Errorf("failed get db: %w", err)
	}

	driverID := sql.NullString{}
	if order.DriverID != nil && *order.DriverID != "" {
		driverID = sql.NullString{String: *order.DriverID, Valid: true}
	}

	originAddr := sql.NullString{}
	if order.OriginAddress != "" {
		originAddr = sql.NullString{String: order.OriginAddress, Valid: true}
	}

	destAddr := sql.NullString{}
	if order.DestinationAddress != "" {
		destAddr = sql.NullString{String: order.DestinationAddress, Valid: true}
	}

	estimatedFare := sql.NullFloat64{}
	if order.EstimatedFare != nil {
		estimatedFare = sql.NullFloat64{Float64: *order.EstimatedFare, Valid: true}
	}

	distanceKm := sql.NullFloat64{}
	if order.DistanceKm != nil {
		distanceKm = sql.NullFloat64{Float64: *order.DistanceKm, Valid: true}
	}

	distanceActual := sql.NullFloat64{}
	if order.DistanceActual != nil {
		distanceActual = sql.NullFloat64{Float64: *order.DistanceActual, Valid: true}
	}

	durationActual := sql.NullString{}
	if order.DurationActual != nil {
		durationActual = sql.NullString{String: *order.DurationActual, Valid: true}
	}

	status := defaultString(order.Status, "REQUESTED")
	paymentMethod := defaultString(order.PaymentMethod, "WALLET")
	paymentStatus := defaultString(order.PaymentStatus, "UNPAID")

	query := `
		INSERT INTO orders (
			order_id,
			passenger_id,
			driver_id,
			origin_lat,
			origin_lng,
			destination_lat,
			destination_lng,
			origin_address,
			destination_address,
			min_price,
			max_price,
			best_route_km,
			best_route_price,
			best_route_duration,
			status,
			payment_method,
			payment_status,
			estimated_fare,
			distance_km,
			distance_actual,
			duration_actual
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	_, err = db.ExecContext(ctx, query,
		order.OrderID,
		order.PassengerID,
		driverID,
		order.OriginLat,
		order.OriginLng,
		order.DestinationLat,
		order.DestinationLng,
		originAddr,
		destAddr,
		order.MinPrice,
		order.MaxPrice,
		order.BestRouteKm,
		order.BestRoutePrice,
		order.BestRouteDuration,
		status,
		paymentMethod,
		paymentStatus,
		estimatedFare,
		distanceKm,
		distanceActual,
		durationActual,
	)

	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, req *entity.UpdateOrderRequest) error {
	db, err := r.DB.GetDB()
	if err != nil {
		return err
	}

	query := `
		UPDATE orders SET
			order_id = ?,
			passenger_id = ?,
			driver_id = ?,
			origin_lat = ?,
			origin_lng = ?,
			destination_lat = ?,
			destination_lng = ?,
			origin_address = ?,
			destination_address = ?,
			min_price = ?,
			max_price = ?,
			best_route_km = ?,
			best_route_price = ?,
			best_route_duration = ?,
			status = ?,
			payment_method = ?,
			payment_status = ?
		WHERE id = ?
	`

	var driverID sql.NullString
	if req.DriverID != nil && *req.DriverID != "" {
		driverID = sql.NullString{String: *req.DriverID, Valid: true}
	}

	_, err = db.ExecContext(ctx, query,
		req.OrderID,
		req.PassengerID,
		driverID,
		req.OriginLat,
		req.OriginLng,
		req.DestinationLat,
		req.DestinationLng,
		req.OriginAddress,
		req.DestinationAddress,
		req.MinPrice,
		req.MaxPrice,
		req.BestRouteKm,
		req.BestRoutePrice,
		req.BestRouteDuration,
		req.Status,
		req.PaymentMethod,
		req.PaymentStatus,
		req.ID,
	)

	return err
}

func (r *OrderRepository) AssignDriverToOrder(ctx context.Context, orderID string, passengerID string, driverID string) (bool, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return false, err
	}

	query := `
		UPDATE orders
		SET driver_id = ?, status = 'ACCEPTED'
		WHERE order_id = ?
		  AND passenger_id = ?
		  AND (driver_id IS NULL OR driver_id = '')
		  AND status IN ('REQUESTED','MATCHING')
	`

	res, err := db.ExecContext(ctx, query, driverID, orderID, passengerID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *OrderRepository) UpdateStatusOrder(ctx context.Context, orderID string, status string) (bool, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return false, err
	}

	query := `
		UPDATE orders
		SET status = ?
		WHERE order_id = ?
	`

	res, err := db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *OrderRepository) UpdateStatusOrderForDriver(ctx context.Context, orderID, driverID, fromStatus, toStatus string) (bool, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return false, err
	}

	query := `
		UPDATE orders
		SET status = ?
		WHERE order_id = ?
		  AND driver_id = ?
		  AND status = ?
	`

	res, err := db.ExecContext(ctx, query, toStatus, orderID, driverID, fromStatus)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *OrderRepository) CompleteTrip(ctx context.Context, orderID, driverID string, distanceActual float64, durationActual string) (bool, error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return false, err
	}

	query := `
		UPDATE orders
		SET 
			status = 'COMPLETED',
			distance_actual = ?,
			duration_actual = ?,
			updated_at = NOW()
		WHERE order_id = ?
		  AND driver_id = ?
		  AND status = 'ON_GOING'
	`

	res, err := db.ExecContext(ctx, query, distanceActual, durationActual, orderID, driverID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func defaultString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
