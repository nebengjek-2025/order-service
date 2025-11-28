package model

type UserEvent struct {
	ID      string      `json:"id,omitempty"`
	Message RequestRide `json:"message,omitempty"`
}

type DriverMatchEvent struct {
	EventID      string       `json:"event_id"`
	OrderID      string       `json:"order_id"`
	PassengerID  string       `json:"passenger_id"`
	DriverID     string       `json:"driver_id"`
	RouteSummary RouteSummary `json:"route_summary,omitempty"`
}

func (u *UserEvent) GetId() string {
	return u.ID
}

func (e *DriverMatchEvent) GetId() string {
	return e.EventID
}
