package usecase

import (
	"context"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int) (*model.Courier, error)
	GetAll(ctx context.Context) ([]model.Courier, error)
	Create(ctx context.Context, courier *model.CreateCourierRequest) (int, error)
	Update(ctx context.Context, courier *model.UpdateCourierRequest) error
	Delete(ctx context.Context, id int) error
}

type DeliveryRepository interface {
	AssignCourierWithUpdate(ctx context.Context, courierId int, orderId string, deadline time.Time) (*model.Delivery, error)
	UnAssign(ctx context.Context, req *model.UnAssignDeliveryRequest) (*model.UnAssignedDeliveryResponse, error) 
}