package handler

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
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.OrderId == "" {
		http.Error(w, `{"error": "Missing orderId"}`,  http.StatusBadRequest)
		return
	}

	response, err := c.useCase.AssignCourier(r.Context(), req)
	if err != nil {
		switch err {
		case model.ErrNoAvailableCouriers:
			http.Error(w, `{"error": "No available couriers"}`, http.StatusConflict)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	c.writeJSON(w, http.StatusOK, response)
}

func (c *DeliveryController) Unassign(w http.ResponseWriter, r *http.Request) {
	var req dto.UnassignDeliveryRequest 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.OrderId == "" {
		http.Error(w, `{"error": "Missing orderId"}`,  http.StatusBadRequest)
		return
	}

	response, err := c.useCase.UnassignCourier(r.Context(), req)
	if err != nil {
		switch err {
		case model.ErrNoRelationFound:
			http.Error(w, `{"error": "Delivery not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	c.writeJSON(w, http.StatusOK, response)
}



func (c *DeliveryController) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}