package model

type OrderDetailRequest struct {
	UserID  string `json:"userId" validate:"required"`
	OrderID string `json:"orderId" validate:"required"`
}

type ConfirmOrderRequest struct {
	UserID   string `json:"userId" validate:"required"`
	DriverID string `json:"driverId" validate:"required"`
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

type RequestRide struct {
	RouteSummary RouteSummary `json:"routeSummary" bson:"routeSummary"`
	UserId       string       `json:"userId" bson:"userId"`
}

type Route struct {
	Origin      LocationRequest `json:"origin" `
	Destination LocationRequest `json:"destination"`
}

type FindDriverResponse struct {
	Message string      `json:"message"`
	Driver  interface{} `json:"driver"`
}

type FindDriverRequest struct {
	UserID        string `json:"userId" validate:"required"`
	PaymentMethod string `json:"paymentMethod" validate:"required,oneof=wallet cash qris"`
}
