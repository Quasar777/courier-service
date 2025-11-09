package api

// TODO: исправить все возвращения ошибок на http.Error(...)

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/Quasar777/courier-service/repository"

	"github.com/go-chi/chi/v5"
)

type CourierController struct {
	repository repository.CourierRepository
}

func NewCourierController(r repository.CourierRepository) *CourierController {
	return &CourierController{repository: r}
}

func (c *CourierController) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	
	var courierDB *model.CourierDB
	courierDB, err = c.repository.GetOneById(ctx, id)

	if err != nil {
		if strings.Contains(err.Error(), "courier with a such id is not found") {
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courierDB)
}


func (c *CourierController) GetMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	

	couriers, err := c.repository.GetAll(ctx) 
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	var reqCourier model.CreateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&reqCourier); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	if reqCourier.Name == "" || reqCourier.Lastname == "" || reqCourier.Phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"fields 'name', 'lastname', 'phone' must be filled",
		})
		return
	}

	ctx := r.Context()

	id, err := c.repository.Create(ctx, &reqCourier)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, `{"error": "courier with this phone is already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": id,
		"message": "Courier created succesfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	var reqCourier model.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&reqCourier); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	if reqCourier.Name == "" || reqCourier.Lastname == "" ||
	   reqCourier.Phone == "" || reqCourier.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"fields 'name', 'lastname', 'phone', 'status' must be filled",
		})
		return
	}

	ctx := r.Context()

	err := c.repository.Update(ctx, &reqCourier)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, `{"error": "Courier with this phone is already exists"}`, http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string {
		"message": "Profile updated successfully",
	}

	json.NewEncoder(w).Encode(response)
}

func (c *CourierController) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}

	err = c.repository.Delete(r.Context(), id)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}