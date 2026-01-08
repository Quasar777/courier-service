package delivery

import (
	"context"

	"github.com/Quasar777/courier-service/internal/handler/dto"
)

type DeliveryUseCase interface {
	AssignCourier(ctx context.Context, req dto.AssignDeliveryRequest) (*dto.AssignedDeliveryResponse, error)
	UnassignCourier(ctx context.Context, req dto.UnassignDeliveryRequest) (*dto.UnassignedDeliveryResponse, error) 
}