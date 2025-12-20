package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func getConnectionString() (string, error) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		return "", err
	}

	connString := os.Getenv("DB_CONNECTION_STRING_TEST")
	if connString == "" {
		return "", err
	}

	return connString, nil
}

func MustInitPool(ctx context.Context) (*pgxpool.Pool, error) {
	connString, err := getConnectionString()
	if err != nil {
		return nil, err
	}

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
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
		return nil, pingErr
	}

	log.Println("Database test connection pool established")
	return pool, nil
}