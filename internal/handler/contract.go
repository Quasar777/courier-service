package handler

import (
	"context"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
)

type CourierUseCase interface {
	GetCourier(ctx context.Context, id int) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, req dto.CreateCourierRequest) (int, error)
	UpdateCourier(ctx context.Context, req dto.UpdateCourierRequest) error
	DeleteCourier(ctx context.Context, id int) error
}

type DeliveryUseCase interface {
	AssignCourier(ctx context.Context, req dto.AssignDeliveryRequest) (*dto.AssignedDeliveryResponse, error)
	UnassignCourier(ctx context.Context, req dto.UnassignDeliveryRequest) (*dto.UnassignedDeliveryResponse, error) 
}