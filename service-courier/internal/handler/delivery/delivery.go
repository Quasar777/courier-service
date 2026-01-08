package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
)

type DeliveryController struct {
	useCase DeliveryUseCase
}

func NewDeliveryController(u DeliveryUseCase) *DeliveryController {
	return &DeliveryController{useCase:  u}
}

func (c *DeliveryController) Assign(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErrorJSON(w, model.ErrInvalidJSON)
		return
	}

	if req.OrderId == "" {
		responseErrorJSON(w, model.ErrMissingRequiredFields)
		return
	}

	response, err := c.useCase.AssignCourier(r.Context(), req)
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (c *DeliveryController) Unassign(w http.ResponseWriter, r *http.Request) {
	var req dto.UnassignDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErrorJSON(w, model.ErrInvalidJSON)
		return
	}

	if req.OrderId == "" {
		responseErrorJSON(w, model.ErrMissingRequiredFields)
		return
	}

	response, err := c.useCase.UnassignCourier(r.Context(), req)
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}