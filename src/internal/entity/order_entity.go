package entity

import (
	"time"
)

type OrderDetail struct {
	ID                 uint64  `db:"order_id"`
	OrderID            string  `db:"order_id"`
	PassengerID        string  `db:"passenger_id"`
	DriverID           *string `db:"driver_id"`
	OriginLat          float64 `db:"origin_lat"`
	OriginLng          float64 `db:"origin_lng"`
	DestinationLat     float64 `db:"destination_lat"`
	DestinationLng     float64 `db:"destination_lng"`
	OriginAddress      string  `db:"origin_address"`
	DestinationAddress string  `db:"destination_address"`

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

	Payment PaymentDetail `json:"payment,omitempty"`

	Promo PromoDetail `json:"promo,omitempty"`
}

type Order struct {
	ID                 uint64    `db:"id"                  json:"id"`
	OrderID            string    `db:"order_id"            json:"order_id"`
	PassengerID        string    `db:"passenger_id"        json:"passenger_id"`
	DriverID           *string   `db:"driver_id"           json:"driver_id,omitempty"`
	OriginLat          float64   `db:"origin_lat"          json:"origin_lat"`
	OriginLng          float64   `db:"origin_lng"          json:"origin_lng"`
	DestinationLat     float64   `db:"destination_lat"     json:"destination_lat"`
	DestinationLng     float64   `db:"destination_lng"     json:"destination_lng"`
	OriginAddress      string    `db:"origin_address"      json:"origin_address,omitempty"`
	DestinationAddress string    `db:"destination_address" json:"destination_address,omitempty"`
	MinPrice           float64   `db:"min_price"           json:"min_price"`
	MaxPrice           float64   `db:"max_price"           json:"max_price"`
	BestRouteKm        float64   `db:"best_route_km"       json:"best_route_km"`
	BestRoutePrice     float64   `db:"best_route_price"    json:"best_route_price"`
	BestRouteDuration  string    `db:"best_route_duration" json:"best_route_duration"`
	Status             string    `db:"status"              json:"status"`
	PaymentMethod      string    `db:"payment_method"      json:"payment_method"`
	PaymentStatus      string    `db:"payment_status"      json:"payment_status"`
	EstimatedFare      *float64  `db:"estimated_fare"      json:"estimated_fare,omitempty"`
	DistanceKm         *float64  `db:"distance_km"         json:"distance_km,omitempty"`
	DistanceActual     *float64  `db:"distance_actual"     json:"distance_actual,omitempty"`
	DurationActual     *string   `db:"duration_actual"     json:"duration_actual,omitempty"`
	CreatedAt          time.Time `db:"created_at"          json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"          json:"updated_at"`
}

type OrderFilter struct {
	ID            *uint64
	OrderID       *string
	PassengerID   *string
	DriverID      *string
	Status        *string
	StatusNot     *string
	StatusIn      []string
	PaymentStatus *string
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

type CreateOrder struct {
	OrderID            string   `json:"order_id"`
	PassengerID        string   `json:"passenger_id"`
	DriverID           *string  `json:"driver_id,omitempty"`
	OriginLat          float64  `json:"origin_lat"`
	OriginLng          float64  `json:"origin_lng"`
	DestinationLat     float64  `json:"destination_lat"`
	DestinationLng     float64  `json:"destination_lng"`
	OriginAddress      string   `json:"origin_address,omitempty"`
	DestinationAddress string   `json:"destination_address,omitempty"`
	MinPrice           float64  `json:"min_price"`
	MaxPrice           float64  `json:"max_price"`
	BestRouteKm        float64  `json:"best_route_km"`
	BestRoutePrice     float64  `json:"best_route_price"`
	BestRouteDuration  string   `json:"best_route_duration"`
	Status             string   `json:"status,omitempty"`
	PaymentMethod      string   `json:"payment_method,omitempty"`
	PaymentStatus      string   `json:"payment_status,omitempty"`
	EstimatedFare      *float64 `json:"estimated_fare,omitempty"`
	DistanceKm         *float64 `json:"distance_km,omitempty"`
	DistanceActual     *float64 `json:"distance_actual,omitempty"`
	DurationActual     *string  `json:"duration_actual,omitempty"`
}

type UpdateOrderRequest struct {
	ID                 uint64
	OrderID            string
	PassengerID        string
	DriverID           *string
	OriginLat          float64
	OriginLng          float64
	DestinationLat     float64
	DestinationLng     float64
	OriginAddress      string
	DestinationAddress string
	MinPrice           float64
	MaxPrice           float64
	BestRouteKm        float64
	BestRoutePrice     float64
	BestRouteDuration  string
	Status             string
	PaymentMethod      string
	PaymentStatus      string
}
