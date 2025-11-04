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
	Name     string
	Lastname string
	Phone    string
	Status   string
}

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