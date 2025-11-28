package model

import (
	"order-service/src/internal/entity"
	"time"
)

type PickupPassanger struct {
	DriverID    string `json:"driverId" validate:"required"`
	PassangerID string `json:"passangerId" bson:"passangerId" validate:"required"`
	OrderID     string `json:"orderId" validate:"required"`
}
type OrderDetailRequest struct {
	UserID  string `json:"userId" validate:"required"`
	OrderID string `json:"orderId" validate:"required"`
}

type ConfirmOrderRequest struct {
	UserID   string `json:"userId" validate:"required"`
	DriverID string `json:"driverId" validate:"required"`
	OrderID  string `json:"orderId" validate:"required"`
}

type CancelOrderRequest struct {
	UserID  string `json:"userId" validate:"required"`
	OrderID string `json:"orderId" validate:"required"`
}

type LocationRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Address   string  `json:"address" validate:"required"`
}

type LocationSuggestionRequest struct {
	CurrentLocation LocationRequest `json:"currentLocation" validate:"required"`
	Destination     LocationRequest `json:"destination" validate:"required"`
	UserID          string          `json:"userId" validate:"required"`
}

type RouteSummary struct {
	Route             Route   `json:"route"`
	MinPrice          float64 `json:"minPrice"`
	MaxPrice          float64 `json:"maxPrice"`
	BestRouteKm       float64 `json:"bestRouteKm"`
	BestRoutePrice    float64 `json:"bestRoutePrice"`
	BestRouteDuration string  `json:"bestRouteDuration"`
	Duration          int     `json:"duration"`
}

type BroadcastPickupPassanger struct {
	RouteSummary RouteSummary `json:"routeSummary" bson:"routeSummary"`
	DriverID     string       `json:"driverId" bson:"driverId"`
	SocketID     string       `json:"socketId" bson:"socketId"`
	PassangerID  string       `json:"passangerId" bson:"passangerId"`
}

type DriverMatch struct {
	DriverID  string  `json:"Name"`
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
	Dist      float64 `json:"Dist"`
	GeoHash   int32   `json:"GeoHash"`
}

type RequestRide struct {
	RouteSummary RouteSummary `json:"routeSummary" bson:"routeSummary"`
	OrderTempID  string       `json:"orderTempId" bson:"orderTempId"`
	UserId       string       `json:"userId" bson:"userId"`
	Attempt      int          `json:"attempt"`
}

type Route struct {
	Origin      LocationRequest `json:"origin" `
	Destination LocationRequest `json:"destination"`
}

type FindDriverResponse struct {
	OrderID string      `json:"orderId"`
	Message string      `json:"message"`
	Driver  interface{} `json:"driver"`
}

type FindDriverRequest struct {
	UserID        string `json:"userId" validate:"required"`
	PaymentMethod string `json:"paymentMethod" validate:"required,oneof=wallet cash qris"`
}

type AvailableDriverResponse struct {
	DriverID       string `json:"driver_id"`
	Status         string `json:"status"`
	LastSeenAt     string `json:"last_seen_at"`
	City           string `json:"city"`
	Province       string `json:"province"`
	JenisKendaraan string `json:"jenis_kendaraan"`
	Nopol          string `json:"nopol"`
}

type DriversRequest struct {
	DriverID       string `json:"driver_id"`
	Name           string `json:"name"`
	MobileNumber   string `json:"mobile_number"`
	City           string `json:"city"`
	Province       string `json:"province"`
	JenisKendaraan string `json:"jenis_kendaraan"`
	Nopol          string `json:"nopol"`
}

type OrderDriverListRequest struct {
	Order   entity.Order     `json:"order"`
	Drivers []DriversRequest `json:"drivers"`
}

type OrderSummary struct {
	OrderID            string    `json:"order_id"`
	Status             string    `json:"status"`
	OriginAddress      string    `json:"origin_address,omitempty"`
	DestinationAddress string    `json:"destination_address,omitempty"`
	BestRouteDuration  string    `json:"best_route_duration,omitempty"`
	BestRoutePrice     float64   `json:"best_route_price,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type DriverPickupInfo struct {
	DriverID    string `json:"driver_id"`
	Name        string `json:"name"`
	Vehicle     string `json:"vehicle"`
	PlateNumber string `json:"plate_number"`
	City        string `json:"city"`
}

type OrderPickupSummaryResponse struct {
	Order   OrderSummary       `json:"order"`
	Drivers []DriverPickupInfo `json:"drivers"`
}

type ConfirmOrderResponse struct {
	OrderID  string `json:"order_id"`
	UserID   string `json:"user_id"`
	DriverID string `json:"driver_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}
