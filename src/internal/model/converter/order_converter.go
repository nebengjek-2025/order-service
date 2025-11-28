package converter

import (
	"order-service/src/internal/model"
)

func OrderToEvent(order *model.PickupPassanger) *model.OrderEvent {
	return &model.OrderEvent{
		Message: *order,
		ID:      order.OrderID,
	}
}
