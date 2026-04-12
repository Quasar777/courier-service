package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/config"
)

func TestMustInitPool_InvalidConnString(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		DBConnectionString: "invalid://dsn",
		DBMaxConns:         1,
		DBMinConns:         0,
		DBMaxConnLifetime:  time.Second,
	}

	pool, err := mustInitPool(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for invalid connection string, got nil")
	}
	if pool != nil {
		t.Fatal("expected nil pool on error")
	}
}

func TestStartServerWithGS_GracefulShutdown(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: http.NewServeMux(),
	}

	done := make(chan struct{})
	go func() {
		startServerWithGS(ctx, srv, 0)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("startServerWithGS did not stop after context cancellation")
	}
}

func TestStartServerWithGS_AlreadyCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: http.NewServeMux(),
	}

	done := make(chan struct{})
	go func() {
		startServerWithGS(ctx, srv, 0)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("startServerWithGS did not return for canceled context")
	}
}
