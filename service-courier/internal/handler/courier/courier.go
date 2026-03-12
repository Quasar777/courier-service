package courier

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
	
	"github.com/go-chi/chi/v5"
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
	var req dto.CreateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErrorJSON(w, model.ErrInvalidJSON)
		return
	}


	input := toUseCaseCreateModel(req)
	id, err := c.useCase.CreateCourier(r.Context(), input)
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

	input := toUseCaseUpdateModel(req)
	err := c.useCase.UpdateCourier(r.Context(), input)

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

