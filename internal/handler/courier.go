package handler

// TODO: добавить application/json в ответах

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/handler/validation"
	"github.com/Quasar777/courier-service/internal/model"

	"github.com/go-chi/chi/v5"
)

type CourierController struct {
	useCase CourierUseCase
}

func NewCourierController(u CourierUseCase) *CourierController {
	return &CourierController{useCase:  u}
}

func (c *CourierController) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}
	
	var courier *model.Courier
	courier, err = c.useCase.GetCourier(r.Context(), id)

	if err != nil {
		switch err {
		case model.ErrCourierNotFound:
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}
	
	c.writeJSON(w, http.StatusOK, courier)
}


func (c *CourierController) GetMany(w http.ResponseWriter, r *http.Request) {
	couriers, err := c.useCase.GetCouriers(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, http.StatusOK, couriers)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	var reqCourier dto.CreateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&reqCourier); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if err := validation.ValidateRequest[dto.CreateCourierRequest](&reqCourier); err != nil {
		responseErrorJSON(w, err)
		return
	}

	id, err := c.useCase.CreateCourier(r.Context(), reqCourier)
	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	response := map[string]interface{}{
		"id": id,
		"message": "Courier created succesfully",
	}

	c.writeJSON(w, http.StatusCreated, response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if req.Id == 0 {
		http.Error(w, `{"error": "Id is required"}`, http.StatusBadRequest)
		return
	}

	if err := validation.ValidateRequest[dto.UpdateCourierRequest](&req); err != nil {
		responseErrorJSON(w, err)
		return
	}

	err := c.useCase.UpdateCourier(r.Context(), req)

	if err != nil {
		responseErrorJSON(w, err)
		return
	}

	response := map[string]string {
		"message": "Courier updated successfully",
	}

	c.writeJSON(w, http.StatusOK, response)
}

func (c *CourierController) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}

	err = c.useCase.DeleteCourier(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string {
		"message": "Courier deleted successfully",
	}

	c.writeJSON(w, http.StatusOK, response)
}

func (c *CourierController) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func responseErrorJSON(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrMissingRequiredFields):
		http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
	case errors.Is(err, model.ErrInvalidStatus):
		http.Error(w, `{"error": "Invalid status"}`, http.StatusBadRequest)
	case errors.Is(err, model.ErrInvalidPhone):
		http.Error(w, `{"error": "Invalid phone number"}`, http.StatusBadRequest)
	case errors.Is(err, model.ErrInvalidCourierTransportType):
		http.Error(w, `{"error": "Invalid courier transport type"}`, http.StatusBadRequest)
	case errors.Is(err, model.ErrPhoneConflict):
		http.Error(w, `{"error": "Courier with this phone already exists"}`, http.StatusConflict)
	case errors.Is(err, model.ErrCourierNotFound):
		http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
	default:
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
	}
}