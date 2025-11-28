package usecase

import (
	"context"
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

	reqUpdateCourierStatus := &model.UpdateCourierRequest{
		Id: result.CourierId,
		Name: assignedCourier.Name,
		Lastname: assignedCourier.Lastname,
		Phone: assignedCourier.Phone,
		Status: "busy",
		TransportType: assignedCourier.TransportType,
	}

	err = u.CourierRepo.Update(ctx, reqUpdateCourierStatus)
	if err != nil {
		return nil, err
	}

	result.TransportType = assignedCourier.TransportType
	result.Deadline = u.DeadlineFactory.Deadline(time.Now(), result.TransportType)

	_, err = u.DeliveryRepo.Assign(ctx, result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *DeliveryUseCase) UnAssignCourier(ctx context.Context, req *model.UnAssignDeliveryRequest) (*model.UnAssignedDeliveryResponse, error) {
	return nil, nil
}


func pickRandomCourierId(couriers []model.Courier) int {
	ids := []int{}
	for _, courier := range couriers {
		ids = append(ids, courier.Id)
	}

	return ids[rand.Intn(len(ids))]
}