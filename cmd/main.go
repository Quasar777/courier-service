package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Quasar777/courier-service/internal/handler"
	"github.com/Quasar777/courier-service/internal/repository"
	"github.com/Quasar777/courier-service/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	defaultPort     = "8080"
	shutdownTimeout = 5 * time.Second
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error when loading .env file:", err)
	}

	port := getEnv("SERVER_PORT", defaultPort)
	flagPort := flag.String("port", port, "specifying a port")
	flag.Parse()
	if *flagPort != port {
		port = *flagPort
	}

	pool, err := mustInitPool(context.Background())
    if err != nil {
        log.Fatal("error connecting to database: ", err)
    }
    defer pool.Close()

	courierRepository := repository.NewCourierRepository(pool)
	courierUseCase := usecase.NewCourierUseCase(courierRepository)
	courier := handler.NewCourierController(courierUseCase)

	deliveryRepository := repository.NewDeliveryRepository(pool)
	deliveryUseCase := usecase.NewDeliveryUseCase(deliveryRepository, courierRepository)
	delivery := handler.NewDeliveryController(deliveryUseCase)
	
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: initRouter(courier, delivery),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	
	go func() {
		log.Println("starting courier-service on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server starting failed: %v", err)
		}
	}() 
	
	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	
	log.Printf("shutting down server gracefully")
	if err = srv.Shutdown(shutDownCtx); err != nil {
		log.Println("error when shutting down:", err)
	} else {
		log.Println("server stopped")
	}
}

func initRouter(courier *handler.CourierController, delivery *handler.DeliveryController) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/courier/{id}", courier.Get)
	r.Get("/couriers", courier.GetMany)
	r.Post("/couriers", courier.Create)
	r.Put("/courier", courier.Update)
	r.Delete("/courier/{id}", courier.Delete)

	r.Post("/delivery/assign", delivery.Assign)
	r.Post("/delivery/unassign", delivery.Unassign)

	return r
}

func mustInitPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour
	cfg.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	pingAttemptsLimit := 3
	var pingErr error

	for i := range pingAttemptsLimit {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = pool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Printf("db ping attempt %d failed: %v", i, pingErr)
		if i < pingAttemptsLimit {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if pingErr != nil {
		log.Fatalf("Unable to ping database")
	}
	
	log.Println("Database connection pool established")
	return pool, nil
}

func getConnectionString() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		log.Fatal("DB_CONNECTION_STRING not set in .env")
	}

	return connString
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}