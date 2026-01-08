package delivery

import (
	"context"
	"fmt"
	"time"

	"github.com/Quasar777/courier-service/internal/handler/dto"
)

type DeliveryUseCase struct {
	DeliveryRepo DeliveryRepository
	CourierRepo CourierRepository
	DeadlineFactory deadlineCalculatorFactory
}

func NewDeliveryUseCase(deliveryRepo DeliveryRepository, courierRepo CourierRepository, deadlineFactory deadlineCalculatorFactory) *DeliveryUseCase {
	return &DeliveryUseCase{
		DeliveryRepo: deliveryRepo,
		CourierRepo: courierRepo,
		DeadlineFactory: deadlineFactory,
	}
}

func (u *DeliveryUseCase) AssignCourier(ctx context.Context, req dto.AssignDeliveryRequest) (*dto.AssignedDeliveryResponse, error) {
	result := dto.AssignedDeliveryResponse{}

	id, err := u.CourierRepo.GetCourierIdWithFewestOrders(ctx)
	if err != nil {
		return nil, err
	}

	result.CourierId = id
	result.OrderId = req.OrderId

	assignedCourier, err := u.CourierRepo.GetOneById(ctx, result.CourierId)
	if err != nil {
		return nil, err
	}

	result.TransportType = string(assignedCourier.TransportType)

	dc := u.DeadlineFactory.GetDeliveryCalculator(assignedCourier.TransportType)
	// result.Deadline = u.DeadlineFactory.Deadline(time.Now(), result.TransportType)
	result.Deadline = dc.CalculateDeadline()

	_, err = u.DeliveryRepo.AssignCourierWithUpdate(ctx, result.CourierId, result.OrderId, result.Deadline)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *DeliveryUseCase) UnassignCourier(ctx context.Context, req dto.UnassignDeliveryRequest) (*dto.UnassignedDeliveryResponse, error) {
	result := dto.UnassignedDeliveryResponse{
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
			err := u.ReleaseExpiredDeliveries(ctx)
			if err != nil {
				fmt.Println("failed to release expired deliveries:", err)
				return
			}
		}
	}
}

func (u *DeliveryUseCase) ReleaseExpiredDeliveries(ctx context.Context) error {
    return u.CourierRepo.ReleaseCouriers(ctx)
}