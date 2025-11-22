package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Quasar777/courier-service/internal/model"
)

type DeliveryController struct {
	useCase DeliveryUseCase
}

func NewDeliveryController(u DeliveryUseCase) *DeliveryController {
	return &DeliveryController{useCase:  u}
}

func (c *DeliveryController) Assign(w http.ResponseWriter, r *http.Request) {
	var req model.AssignDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON}`, http.StatusBadRequest)
		return
	}

	c.writeJSON(w, 228, map[string]string{"message": "pong from assign"})
}

func (c *DeliveryController) Unassign(w http.ResponseWriter, r *http.Request) {
	var req model.UnAssignDeliveryRequest 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON}`, http.StatusBadRequest)
		return
	}

	c.writeJSON(w, 228, map[string]string{"message": "pong from unassign"})
}



func (c *DeliveryController) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}