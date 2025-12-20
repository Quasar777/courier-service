package delivery

import (
	"context"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int) (*model.Courier, error)
	ReleaseCouriers(ctx context.Context) error
	GetCourierIdWithFewestOrders(ctx context.Context) (int, error)
}

type DeliveryRepository interface {
	AssignCourierWithUpdate(ctx context.Context, courierId int, orderId string, deadline time.Time) (*model.Delivery, error)
	UnassignWithUpdate(ctx context.Context, orderId string) (*model.Delivery, error)
}

type IDeadlineFactory interface {
	Deadline(now time.Time, transportType string) time.Time
}