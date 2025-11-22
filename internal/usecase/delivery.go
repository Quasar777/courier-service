package usecase

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
)

type DeliveryUseCase struct {
	repository DeliveryRepository
}

func NewDeliveryUseCase(r DeliveryRepository) *DeliveryUseCase {
	return &DeliveryUseCase{repository: r}
}

func (u *DeliveryUseCase) AssignCourier(ctx context.Context, orderId model.AssignDeliveryRequest) (model.AssignedDeliveryResponse, error) {
	return model.AssignedDeliveryResponse{}, nil
}

func (u *DeliveryUseCase) 	UnAssignCourier(ctx context.Context, orderId model.UnAssignDeliveryRequest) (model.UnAssignedDeliveryResponse, error) {
	return model.UnAssignedDeliveryResponse{}, nil
}
