package delivery

import (
	"github.com/Quasar777/courier-service/internal/handler/dto"
	uc "github.com/Quasar777/courier-service/internal/usecase/delivery"
)

func toAssignInput(req dto.AssignDeliveryRequest) uc.AssignInput {
	return uc.AssignInput{
		OrderID: req.OrderId,
	}
}

func toAssignResponse(out *uc.AssignOutput) dto.AssignedDeliveryResponse {
	return dto.AssignedDeliveryResponse{
		CourierId: out.CourierId,
		OrderId: out.OrderId,
		TransportType: out.TransportType,
		Deadline: out.Deadline,
	}
}

func toUnassignInput(req dto.UnassignDeliveryRequest) uc.UnassignInput {
	return uc.UnassignInput{
		OrderID: req.OrderId,
	}
}

func toUnassignResponse(out *uc.UnassignOutput) dto.UnassignedDeliveryResponse {
	return dto.UnassignedDeliveryResponse{
		OrderId: out.OrderId,
    	CourierId: out.CourierId,
    	Status: out.Status,
	}
}