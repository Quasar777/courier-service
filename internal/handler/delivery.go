package handler

type DeliveryController struct {
	useCase DeliveryUseCase
}

func NewDeliveryController(u DeliveryUseCase) *DeliveryController {
	return &DeliveryController{useCase:  u}
}