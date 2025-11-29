package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
)

type DeliveryUseCase struct {
	DeliveryRepo DeliveryRepository
	CourierRepo CourierRepository
	DeadlineFactory DeadlineFactory
}

func NewDeliveryUseCase(deliveryRepo DeliveryRepository, courierRepo CourierRepository, deadlineFactory DeadlineFactory) *DeliveryUseCase {
	return &DeliveryUseCase{
		DeliveryRepo: deliveryRepo,
		CourierRepo: courierRepo,
		DeadlineFactory: deadlineFactory,
	}
}

func (u *DeliveryUseCase) AssignCourier(ctx context.Context, req model.AssignDeliveryRequest) (*model.AssignedDeliveryResponse, error) {
	result := model.AssignedDeliveryResponse{}

	allCouriers, err := u.CourierRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	availableCouriers := make([]model.Courier, 0)
	for _, courier := range allCouriers {
		if courier.Status == "available" {
			availableCouriers = append(availableCouriers, courier)
		}
	}

	if len(availableCouriers) == 0 {
		return nil, model.ErrNoAvailableCouriers
	}

	result.CourierId = pickRandomCourierId(availableCouriers)
	result.OrderId = req.OrderId

	assignedCourier, err := u.CourierRepo.GetOneById(ctx, result.CourierId)
	if err != nil {
		return nil, err
	}

	result.TransportType = assignedCourier.TransportType
	result.Deadline = u.DeadlineFactory.Deadline(time.Now(), result.TransportType)

	_, err = u.DeliveryRepo.AssignCourierWithUpdate(ctx, result.CourierId, result.OrderId, result.Deadline)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *DeliveryUseCase) UnassignCourier(ctx context.Context, req model.UnassignDeliveryRequest) (*model.UnassignedDeliveryResponse, error) {
	result := model.UnassignedDeliveryResponse{
		OrderId: req.OrderId,
	}

	d, err := u.DeliveryRepo.UnassignWithUpdate(ctx, req.OrderId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	

	result.CourierId = d.CourierId
	result.Status = "unassigned"

	return &result, nil
}


func pickRandomCourierId(couriers []model.Courier) int {
	ids := []int{}
	for _, courier := range couriers {
		ids = append(ids, courier.Id)
	}

	return ids[rand.Intn(len(ids))]
}