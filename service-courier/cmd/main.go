package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Quasar777/courier-service/internal/config"
	courierHandler "github.com/Quasar777/courier-service/internal/handler/courier"
	deliveryHandler "github.com/Quasar777/courier-service/internal/handler/delivery"
	courierRepo "github.com/Quasar777/courier-service/internal/repository/courier"
	deliveryRepo "github.com/Quasar777/courier-service/internal/repository/delivery"
	"github.com/Quasar777/courier-service/internal/usecase/courier"
	"github.com/Quasar777/courier-service/internal/usecase/delivery"
	worker "github.com/Quasar777/courier-service/internal/usecase/delivery/worker"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const shutdownTimeout = 5 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	flagPort := flag.Int("port", cfg.ServerPort, "specifying a port")
	flag.Parse()
	cfg.ServerPort = *flagPort

	dbPool, err := mustInitPool(context.Background(), cfg)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	defer dbPool.Close()

	deadlineFactory := delivery.NewFactory()

	// Repositories
	courierRepository := courierRepo.NewCourierRepository(dbPool)
	deliveryRepository := deliveryRepo.NewDeliveryRepository(dbPool)

	// Use Cases
	courierUseCase := courier.NewCourierUseCase(courierRepository)
	deliveryUseCase := delivery.NewDeliveryUseCase(deliveryRepository, courierRepository, deadlineFactory)
	deliveryWorker := worker.NewWorker(courierRepository)

	// HTTP Handlers
	courierController := courierHandler.NewCourierController(courierUseCase)
	deliveryController := deliveryHandler.NewDeliveryController(deliveryUseCase)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: initRouter(courierController, deliveryController),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go startDeliveryChecker(ctx, deliveryWorker, cfg.DeliveryCheckerInterval)

	startServerWithGS(ctx, srv, cfg.ServerPort)
}

func initRouter(courier *courierHandler.CourierController, delivery *deliveryHandler.DeliveryController) *chi.Mux {
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

func mustInitPool(ctx context.Context, appCfg config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(appCfg.DBConnectionString)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = appCfg.DBMaxConns
	poolCfg.MaxConnLifetime = appCfg.DBMaxConnLifetime
	poolCfg.MinConns = appCfg.DBMinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	pingAttemptsLimit := 3
	var pingErr error

	for i := 0; i < pingAttemptsLimit; i++ {
		pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
		pingErr = pool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			break
		}

		log.Printf("db ping attempt %d failed: %v", i+1, pingErr)
		if i < pingAttemptsLimit-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if pingErr != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", pingErr)
	}

	log.Println("Database connection pool established")
	return pool, nil
}

func startDeliveryChecker(ctx context.Context, worker *worker.Worker, interval time.Duration) {
	worker.RunDeliveryChecker(ctx, interval)
}

func startServerWithGS(ctx context.Context, srv *http.Server, port int) {
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
	if err := srv.Shutdown(shutDownCtx); err != nil {
		log.Println("error when shutting down:", err)
	} else {
		log.Println("server stopped")
	}
}
