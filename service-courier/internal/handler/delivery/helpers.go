package delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Quasar777/courier-service/internal/model"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	jsonErrNoAvailableCouriers = ErrorResponse{Error: "No available couriers"}
	jsonErrDeliveryNotFound    = ErrorResponse{Error: "Delivery not found"}
	jsonErrInternalServer      = ErrorResponse{Error: "Internal server error"}
	jsonErrInvalidJSON         = ErrorResponse{Error: "Invalid JSON"}
	jsonErrMissingFields       = ErrorResponse{Error: "Missing required fields"}
	jsonErrInvalidOrderID      = ErrorResponse{Error: "Invalid order id"}
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func responseErrorJSON(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidJSON):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidJSON)

	case errors.Is(err, model.ErrMissingRequiredFields):
		writeJSON(w, http.StatusBadRequest, jsonErrMissingFields)

	case errors.Is(err, model.ErrNoAvailableCouriers):
		writeJSON(w, http.StatusConflict, jsonErrNoAvailableCouriers)

	case errors.Is(err, model.ErrNoRelationFound):
		writeJSON(w, http.StatusNotFound, jsonErrDeliveryNotFound)
	
	case errors.Is(err, model.ErrInvalidOrderID):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidOrderID)

	case errors.Is(err, model.ErrInternal):
		writeJSON(w, http.StatusInternalServerError, jsonErrInternalServer)

	default:
		writeJSON(w, http.StatusInternalServerError, jsonErrInternalServer)
	}
}
