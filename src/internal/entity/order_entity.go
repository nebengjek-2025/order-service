package entity

import "time"

type OrderDetail struct {
	OrderID            uint64  `db:"order_id"`
	PassengerID        uint64  `db:"passenger_id"`
	DriverID           *uint64 `db:"driver_id"`
	OriginLat          float64 `db:"origin_lat"`
	OriginLng          float64 `db:"origin_lng"`
	DestinationLat     float64 `db:"destination_lat"`
	DestinationLng     float64 `db:"destination_lng"`
	OriginAddress      string  `db:"origin_address"`
	DestinationAddress string  `db:"destination_address"`

	Route             []byte  `db:"route"` // JSON
	MinPrice          float64 `db:"min_price"`
	MaxPrice          float64 `db:"max_price"`
	BestRouteKm       float64 `db:"best_route_km"`
	BestRoutePrice    float64 `db:"best_route_price"`
	BestRouteDuration string  `db:"best_route_duration"`

	Status        string    `db:"status"`
	PaymentMethod string    `db:"payment_method"`
	PaymentStatus string    `db:"payment_status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`

	// Payment
	Payment PaymentDetail `json:"payment"`

	// Promo
	Promo PromoDetail `json:"promo"`
}

type PaymentDetail struct {
	ID          *uint64    `db:"payment_id"`
	Amount      *float64   `db:"payment_amount"`
	Status      *string    `db:"payment_status_detail"`
	Provider    *string    `db:"payment_provider"`
	ReferenceID *string    `db:"payment_ref_id"`
	PaidAt      *time.Time `db:"payment_paid_at"`
}

type PromoDetail struct {
	RedemptionID  *uint64  `db:"redemption_id"`
	Discount      *float64 `db:"discount_applied"`
	PromoCode     *string  `db:"promo_code"`
	PromoName     *string  `db:"promo_name"`
	DiscountType  *string  `db:"discount_type"`
	DiscountValue *float64 `db:"discount_value"`
	MaxDiscount   *float64 `db:"max_discount"`
}
