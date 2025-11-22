package usecase

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
)

type DeliveryUseCase struct {
	DeliveryRepo DeliveryRepository
	CourierRepo CourierRepository
}

func NewDeliveryUseCase(deliveryRepo DeliveryRepository, courierRepo CourierRepository) *DeliveryUseCase {
	return &DeliveryUseCase{
		DeliveryRepo: deliveryRepo,
		CourierRepo: courierRepo,
	}
}

func (u *DeliveryUseCase) AssignCourier(ctx context.Context, orderId model.AssignDeliveryRequest) (model.AssignedDeliveryResponse, error) {
	


	return model.AssignedDeliveryResponse{}, nil
}

func (u *DeliveryUseCase) 	UnAssignCourier(ctx context.Context, orderId model.UnAssignDeliveryRequest) (model.UnAssignedDeliveryResponse, error) {
	return model.UnAssignedDeliveryResponse{}, nil
}

