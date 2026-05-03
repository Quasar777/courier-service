package orderbus

import "github.com/Quasar777/courier-service/internal/domain/order"

func dtoToDomainOrderID(dto *message) order.OrderID {
	return order.OrderID{OrderID: dto.OrderID}
}
