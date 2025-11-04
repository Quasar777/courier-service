package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Quasar777/courier-service/database"
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

func GetCourier(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	conn := database.InitConnection(ctx)
	defer conn.Close(ctx)

	var courier courier
	err = conn.QueryRow(ctx,
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


func GetCouriers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	conn := database.InitConnection(ctx)
	defer conn.Close(ctx)

	var couriers []courier
	rows, err := conn.Query(ctx,
		"SELECT id, name, lastname, phone, status FROM couriers",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in courier handler:", err)
		return
	}

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

func CreateCourier(w http.ResponseWriter, r *http.Request) {
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

	ctx := context.Background()
	conn := database.InitConnection(ctx)
	defer conn.Close(ctx)

	// check constraints
	var candidatID int
	err := conn.QueryRow(ctx, 
		"SELECT * FROM couriers WHERE phone = $1",
		resCourier.Phone).Scan(&candidatID)
	if err != pgx.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string {
			"error":"corier with this phone number is already exists",
		})
		return
	}

	// create courier
	err = conn.QueryRow(ctx,
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