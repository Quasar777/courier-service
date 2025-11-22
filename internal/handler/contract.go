package handler

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
)

type CourierUseCase interface {
	GetCourier(ctx context.Context, id int) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, req model.CreateCourierRequest) (int, error)
	UpdateCourier(ctx context.Context, req model.UpdateCourierRequest) error
	DeleteCourier(ctx context.Context, id int) error
}

type DeliveryUseCase interface {
	GetDeliveries(ctx context.Context) ([]model.Delivery, error)
	AssignCourier(ctx context.Context, orderId string) (model.AssignedDeliveryResponse, error)
	UnAssignCourier(ctx context.Context, orderId string) (model.UnAssignedDeliveryResponse, error) 
}