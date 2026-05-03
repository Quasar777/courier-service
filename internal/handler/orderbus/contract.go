package orderbus

import (
	"context"

	"github.com/Quasar777/courier-service/internal/domain/order"
	"github.com/Quasar777/courier-service/internal/service/deliveryapp"
)

type eventStrategyFactory interface {
	GetEventStrategy(statusMsg string, statusNow string) (deliveryapp.Executor, error)
}

type orderGateway interface {
	GetOrderByID(ctx context.Context, orderID order.OrderID) (*order.Order, error)
}
