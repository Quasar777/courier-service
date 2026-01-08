package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/integration/testkit"
	courierRepo "github.com/Quasar777/courier-service/internal/repository/courier"
	courierUC "github.com/Quasar777/courier-service/internal/usecase/courier"
	courierHandler "github.com/Quasar777/courier-service/internal/handler/courier"

	"github.com/go-chi/chi/v5"
)

type gotCourier struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Lastname      string `json:"lastname"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transportType"`
}

var pg *testkit.TestPostgres

func TestMain(m *testing.M) {
	// DEBUG
	dir, _ := os.Getwd()
	fmt.Fprintln(os.Stderr, dir)
	fmt.Println("DIR:", dir)

	// Setup
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Minute)
	var err error
	pg, err = testkit.StartPostgres(ctx)
	cancel()
	if err != nil {
		fmt.Fprintln(os.Stderr, "start postgres:", err)
		os.Exit(1)
	}

    if err := testkit.RunMigrations(pg.DB, "../../migrations"); err != nil {
		fmt.Fprintln(os.Stderr, "migrations:", err)
		teardown()
		os.Exit(1)
	}

	// Run
	code := m.Run()

	// Teardown
	teardown()

	os.Exit(code)
}

func teardown() {
	if pg == nil {
		return
	}
	tdCtx, tdCancel := context.WithTimeout(context.Background(), 30*time.Second)
	pg.Close(tdCtx)
	tdCancel()
}

func TestGetCourier_HappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	t.Cleanup(func() {
		_, err := pg.DB.Exec(ctx, `TRUNCATE couriers RESTART IDENTITY CASCADE`)
		if err != nil {
			t.Errorf("failed to clean mock courier: %v", err)
		}
	})

	repo := courierRepo.NewCourierRepository(pg.DB)
	uc := courierUC.NewCourierUseCase(repo)
	h := courierHandler.NewCourierController(uc)

	createReq := &dto.CreateCourierRequest{
		Name: "Ivan",
		Lastname: "Ivanov",
		Phone: "+79991119911",
		TransportType: "on_foot",
	}

	id, err := addCourierToDB(t, ctx, createReq)
	if err != nil {
		t.Fatal("failed to craete mock courier")
	}

	r := chi.NewRouter()
	r.Route("/couriers", func(r chi.Router) {
		r.Get("/{id}", h.Get)
	})

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	idAsStr := strconv.Itoa(id)
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/couriers/"+idAsStr, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var got gotCourier
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type: application/json, got: %v", ct)
	}

	assertCourier(t, got, id, createReq.Name, createReq.Lastname, createReq.Phone, "available", "on_foot")
}

func TestGetCourier_InvalidID(t *testing.T) {
	repo := courierRepo.NewCourierRepository(pg.DB)
	uc := courierUC.NewCourierUseCase(repo)
	h := courierHandler.NewCourierController(uc)

	r := chi.NewRouter()
	r.Route("/couriers", func(r chi.Router) {
		r.Get("/{id}", h.Get)
	})

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	id := "badID"
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/couriers/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type: application/json, got: %v", ct)
	}
}

func TestGetCourier_NotFound(t *testing.T) {
	repo := courierRepo.NewCourierRepository(pg.DB)
	uc := courierUC.NewCourierUseCase(repo)
	h := courierHandler.NewCourierController(uc)

	r := chi.NewRouter()
	r.Route("/couriers", func(r chi.Router) {
		r.Get("/{id}", h.Get)
	})

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	id := "9999999"
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/couriers/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type: application/json, got: %v", ct)
	}
}

func TestGetAllCouriers_HappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	t.Cleanup(func() {
		_, err := pg.DB.Exec(ctx, `TRUNCATE couriers RESTART IDENTITY CASCADE`)
		if err != nil {
			t.Errorf("failed to clean mock courier: %v", err)
		}
	})

	repo := courierRepo.NewCourierRepository(pg.DB)
	uc := courierUC.NewCourierUseCase(repo)
	h := courierHandler.NewCourierController(uc)

	createReq1 := &dto.CreateCourierRequest{
		Name: "Ivan",
		Lastname: "Ivanov",
		Phone: "+79991119911",
		TransportType: "on_foot",
	}
	createReq2 := &dto.CreateCourierRequest{
		Name: "Jim",
		Lastname: "Beem",
		Phone: "+79881118811",
		TransportType: "scooter",
	}

	id1, err := addCourierToDB(t, ctx, createReq1)
	if err != nil {
		t.Fatal("failed to craete mock courier")
	}
	
	id2, err := addCourierToDB(t, ctx, createReq2)
	if err != nil {
		t.Fatal("failed to craete mock courier")
	}

	r := chi.NewRouter()
	r.Route("/couriers", func(r chi.Router) {
		r.Get("/", h.GetMany)
	})

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/couriers", nil)
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var got []gotCourier
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type: application/json, got: %v", ct)
	}

	if len(got) != 2 {
		t.Fatalf("expected len of couriers array %v, got: %v", 2, len(got))
	}

	got1 := got[0]
	got2 := got[1]

	assertCourier(t, got1, id1, createReq1.Name, createReq1.Lastname, createReq1.Phone, "available",  createReq1.TransportType)
	assertCourier(t, got2, id2, createReq2.Name, createReq2.Lastname, createReq2.Phone, "available",  createReq2.TransportType)
}

func assertCourier(t *testing.T, got gotCourier, id int, name, lastname, phone, status, transport string) {
	t.Helper()

	if got.ID != id {
		t.Fatalf("expected id=%v, got %v", id, got.ID)
	}
	if got.Name != name {
		t.Fatalf("expected name=%v, got %v", name, got.Name)
	}
	if got.Lastname != lastname {
		t.Fatalf("expected lastname=%v, got %v", lastname, got.Lastname)
	}
	if got.Phone != phone {
		t.Fatalf("expected phone=%v, got %v", phone, got.Phone)
	}
	if got.Status != status {
		t.Fatalf("expected status=%v, got %v", status, got.Status)
	}
	if got.TransportType != transport {
		t.Fatalf("expected transportType=%v, got %v", transport, got.TransportType)
	}
}

func addCourierToDB(t *testing.T, ctx context.Context, data *dto.CreateCourierRequest) (int, error) {
	t.Helper()

	if data.Status == "" {
		data.Status = "available"
	}

	if data.TransportType == "" {
		data.TransportType = "on_foot"
	}

	var id int
	err := pg.DB.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, data.Name, data.Lastname, data.Phone, data.Status, data.TransportType).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}