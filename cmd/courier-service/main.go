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

	"github.com/Quasar777/courier-service/api"
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

	port := getEnv("PORT", defaultPort)
	flagPort := flag.String("port", port, "specifying a port")
	flag.Parse()
	if *flagPort != port {
		port = *flagPort
	}

	// setting up a router 
	r := chi.NewRouter()

	r.Get("/ping", api.HandlePing)
	r.Head("/healthcheck", api.HandleHealthCheck)

	// setting up a server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: r,
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

	// gracefult shutdown
	shutDownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	
	log.Printf("shutting down server gracefully")
	if err = srv.Shutdown(shutDownCtx); err != nil {
		log.Println("error when shutting down:", err)
	} else {
		log.Println("server stopped")
	}
}

// getEnv возвращает значение env или дефолт
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}