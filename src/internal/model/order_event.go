package model

type OrderEvent struct {
	ID      string          `json:"id,omitempty"`
	Message PickupPassanger `json:"message,omitempty"`
}

func (u *OrderEvent) GetId() string {
	return u.ID
}
