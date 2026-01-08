package courier

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"

	"github.com/go-chi/chi/v5"
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

type CourierController struct {
	useCase CourierUseCase
}

func NewCourierController(u CourierUseCase) *CourierController {
	return &CourierController{useCase: u}
}

func (c *CourierController) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responseErrorJSON(w, model.ErrInvalidId)
		return
	}

	var courier *model.Courier
	courier, err = c.useCase.GetCourier(r.Context(), id)

	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, courier)
}

func (c *CourierController) GetMany(w http.ResponseWriter, r *http.Request) {
	couriers, err := c.useCase.GetCouriers(r.Context())
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, couriers)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	var reqCourier dto.CreateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&reqCourier); err != nil {
		responseErrorJSON(w, model.ErrInvalidJSON)
		return
	}

	id, err := c.useCase.CreateCourier(r.Context(), reqCourier)
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Courier created succesfully",
	}

	writeJSON(w, http.StatusCreated, response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErrorJSON(w, model.ErrInvalidJSON)
		return
	}

	err := c.useCase.UpdateCourier(r.Context(), req)

	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	response := map[string]string{
		"message": "Courier updated successfully",
	}

	writeJSON(w, http.StatusOK, response)
}

func (c *CourierController) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	err = c.useCase.DeleteCourier(r.Context(), id)
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	response := map[string]string{
		"message": "Courier deleted successfully",
	}

	writeJSON(w, http.StatusOK, response)
}

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
