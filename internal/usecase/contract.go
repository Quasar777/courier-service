package usecase

import (
	"context"
	"time"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int) (*model.Courier, error)
	GetAll(ctx context.Context) ([]model.Courier, error)
	Create(ctx context.Context, courier *dto.CreateCourierRequest) (int, error)
	Update(ctx context.Context, courier *dto.UpdateCourierRequest) error
	Delete(ctx context.Context, id int) error
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