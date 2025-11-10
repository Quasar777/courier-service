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

	"github.com/Quasar777/courier-service/internal/database"
	"github.com/Quasar777/courier-service/internal/handlers"
	"github.com/Quasar777/courier-service/internal/repository"
	"github.com/Quasar777/courier-service/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const (
	defaultPort     = "8080"
	shutdownTimeout = 5 * time.Second
)

func main() {
	// env file loading and flags parsing
	err := godotenv.Load()
	if err != nil {
		log.Println("error when loading .env file:", err)
	}

	// Flag parsing
	port := getEnv("SERVER_PORT", defaultPort)
	flagPort := flag.String("port", port, "specifying a port")
	flag.Parse()
	if *flagPort != port {
		port = *flagPort
	}

	// Connection pool init
	pool, err := database.InitPool(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

	courierRepository := *repository.NewCourierRepository(pool)
	courierUseCase := *usecase.NewCourierUseCase(courierRepository)
	courier := handlers.NewCourierController(courierUseCase)
	
	// Setup http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: initRouter(courier),
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

	// Gracefult shutdown
	shutDownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	
	log.Printf("shutting down server gracefully")
	if err = srv.Shutdown(shutDownCtx); err != nil {
		log.Println("error when shutting down:", err)
	} else {
		log.Println("server stopped")
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func initRouter(courier *handlers.CourierController) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/courier/{id}", courier.Get)
	r.Get("/couriers", courier.GetMany)
	r.Post("/couriers", courier.Create)
	r.Put("/courier", courier.Update)
	r.Delete("/courier/{id}", courier.Delete)

	return r
}