package courier

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
	jsonErrNotFound             = ErrorResponse{Error: "Courier not found"}
	jsonErrInternalServer       = ErrorResponse{Error: "Internal server error"}
	jsonErrInvalidId            = ErrorResponse{Error: "Invalid id"}
	jsonErrInvalidJSON          = ErrorResponse{Error: "Invalid JSON"}
	jsonErrMissingFields        = ErrorResponse{Error: "Missing required fields"}
	jsonErrInvalidStatus        = ErrorResponse{Error: "Invalid status"}
	jsonErrInvalidPhone         = ErrorResponse{Error: "Invalid phone number"}
	jsonErrInvalidTransportType = ErrorResponse{Error: "Invalid courier transport type"}
	jsonErrPhoneConflict        = ErrorResponse{Error: "Courier with this phone already exists"}
	jsonErrInvalidName          = ErrorResponse{Error: "Invalid courier name"}
	jsonErrInvalidLastname      = ErrorResponse{Error: "Invalid courier lastname"}
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

	case errors.Is(err, model.ErrInvalidId):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidId)

	case errors.Is(err, model.ErrInvalidStatus):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidStatus)

	case errors.Is(err, model.ErrInvalidPhone):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidPhone)

	case errors.Is(err, model.ErrInvalidCourierName):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidName)

	case errors.Is(err, model.ErrInvalidCourierLastname):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidLastname)

	case errors.Is(err, model.ErrInvalidCourierTransportType):
		writeJSON(w, http.StatusBadRequest, jsonErrInvalidTransportType)

	case errors.Is(err, model.ErrPhoneConflict):
		writeJSON(w, http.StatusConflict, jsonErrPhoneConflict)

	case errors.Is(err, model.ErrCourierNotFound):
		writeJSON(w, http.StatusNotFound, jsonErrNotFound)

	case errors.Is(err, model.ErrInternal):
		writeJSON(w, http.StatusInternalServerError, jsonErrInternalServer)

	default:
		writeJSON(w, http.StatusInternalServerError, jsonErrInternalServer)
	}
}