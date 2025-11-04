package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type courier struct {
	Id       int
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

const (
	defaultStatus = "paused"
)

func (s *Server) GetCourier(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	
	var courier courier
	err = s.DB.QueryRow(ctx,
		"SELECT id, name, lastname, phone, status FROM couriers WHERE id = $1",
		id).
		Scan(&courier.Id, &courier.Name, &courier.Lastname, &courier.Phone, &courier.Status)
	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	
	json.NewEncoder(w).Encode(courier)
}


func (s *Server) GetCouriers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	

	var couriers []courier
	rows, err := s.DB.Query(ctx,
		"SELECT id, name, lastname, phone, status FROM couriers",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var courier courier
		err := rows.Scan(&courier.Id, &courier.Name, &courier.Lastname, &courier.Phone, &courier.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("error in courier handler:", err)
			return
		}
		couriers = append(couriers, courier)
	}

	json.NewEncoder(w).Encode(couriers)
}

func (s *Server) CreateCourier(w http.ResponseWriter, r *http.Request) {
	var resCourier courier
	if err := json.NewDecoder(r.Body).Decode(&resCourier); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	if resCourier.Name == "" || resCourier.Lastname == "" || resCourier.Phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"fields 'name', 'lastname', 'phone' must be filled",
		})
		return
	}

	ctx := r.Context()

	// Check confilcts
	var candidatID int
	err := s.DB.QueryRow(ctx, 
		"SELECT id FROM couriers WHERE phone = $1",
		resCourier.Phone).Scan(&candidatID)
	if err != pgx.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"courier with this phone number is already exists",
		})
		return
	}

	// Create courier
	err = s.DB.QueryRow(ctx,
		"INSERT INTO couriers (name, lastname, phone, status) VALUES ($1, $2, $3, $4) RETURNING id;",
		resCourier.Name, resCourier.Lastname, resCourier.Phone, defaultStatus,
	).Scan(&resCourier.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	
	resCourier.Status = defaultStatus

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resCourier)
}

func (s *Server) UpdateCourier(w http.ResponseWriter, r *http.Request) {
	var resCourier courier
	if err := json.NewDecoder(r.Body).Decode(&resCourier); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}
	if resCourier.Name == "" || resCourier.Lastname == "" ||
	   resCourier.Phone == "" || resCourier.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"fields 'name', 'lastname', 'phone', 'status' must be filled",
		})
		return
	}

	ctx := r.Context()

	var candidatCourier courier
	err := s.DB.QueryRow(ctx, 
		"SELECT id, name, lastname, phone, status FROM couriers WHERE id = $1",
		resCourier.Id).Scan(
			&candidatCourier.Id,
			&candidatCourier.Name,
			&candidatCourier.Lastname,
			&candidatCourier.Phone,
			&candidatCourier.Status,
		)	
	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}

	// Check confilcts (phone)
	var samePhoneCandidatID int
	err = s.DB.QueryRow(ctx, 
		"SELECT id FROM couriers WHERE phone = $1 AND id != $2;",
		resCourier.Phone, resCourier.Id).Scan(&samePhoneCandidatID)
	if err != pgx.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"courier with this phone number is already exists",
		})
		return
	}
	
	// Update
	_, err = s.DB.Exec(ctx,
		"UPDATE couriers SET name = $1, lastname = $2, phone = $3, status = $4 WHERE id = $5",
		resCourier.Name, resCourier.Lastname, resCourier.Phone, resCourier.Status, resCourier.Id,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}

	json.NewEncoder(w).Encode(resCourier)
}