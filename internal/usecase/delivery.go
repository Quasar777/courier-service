package usecase

import (
	"context"
	"fmt"
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

	id, err := u.DeliveryRepo.GetCourierIdWithFewestOrders(ctx)
	if err != nil {
		return nil, err
	}

	result.CourierId = id
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

func (u *DeliveryUseCase) RunDeliveryChecker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopping ticker")
			return
		case <-ticker.C:
			fmt.Println("test tick")
			err := u.ReleaseExpiredDeliveries(ctx)
			if err != nil {
				fmt.Println("failed to release expired deliveries:", err)
				return
			}
		}
	}
}

func (u *DeliveryUseCase) ReleaseExpiredDeliveries(ctx context.Context) error {
    return u.DeliveryRepo.ReleaseCouriers(ctx)
}