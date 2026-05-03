package deliveryhttp

import (
	"github.com/Quasar777/courier-service/internal/domain/delivery"
	"github.com/Quasar777/courier-service/internal/domain/order"
)

func toDomainOrderID(req DeliveryOrderRequest) order.OrderID {
	return order.OrderID{OrderID: req.OrderID}
}

func domainToDTOAssign(del *delivery.AssignResult) DeliveryAssignResponse {
	return DeliveryAssignResponse{
		CourierID:     del.CourierID,
		OrderID:       del.OrderID,
		TransportType: del.TransportType,
		Deadline:      del.Deadline,
	}
}

func domainToDTOUnassign(del *delivery.UnassignResult) DeliveryUnassignResponse {
	return DeliveryUnassignResponse{
		OrderID:   del.OrderID,
		Status:    del.Status,
		CourierID: del.CourierID,
	}
}
