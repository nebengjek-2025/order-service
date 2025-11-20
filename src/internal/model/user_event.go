package model

type UserEvent struct {
	ID      string      `json:"id,omitempty"`
	Message RequestRide `json:"message,omitempty"`
}

func (u *UserEvent) GetId() string {
	return u.ID
}
