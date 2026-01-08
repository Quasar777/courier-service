package delivery

import (
	"context"
	"fmt"
	
	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
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
	if !isOrderIDValid(req.OrderId) {
		return nil, model.ErrInvalidOrderID
	}
	
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
	
	result.Deadline = dc.CalculateDeadline()

	_, err = u.DeliveryRepo.AssignCourierWithUpdate(ctx, result.CourierId, result.OrderId, result.Deadline)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *DeliveryUseCase) UnassignCourier(ctx context.Context, req dto.UnassignDeliveryRequest) (*dto.UnassignedDeliveryResponse, error) {
	if !isOrderIDValid(req.OrderId) {
		return nil, model.ErrInvalidOrderID
	}

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

// Order ID должен быть в формате "XXXXXYYYYY".
// "X" - буквы латинского алфавита, "Y" - цифры
func isOrderIDValid(id string) bool {
	if len(id) != 10 {
		return false
	}

	for i := 0; i < 5; i++ {
		if id[i] < 'A' || (id[i] > 'Z' && id[i] < 'a') || id[i] > 'z' {
			return false
		}
	}

	for i := 5; i < 10; i++ {
		if id[i] < '0' || id[i] > '9' {
			return false
		}
	}

	return true
}