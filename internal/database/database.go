package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitConnection(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, getConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func InitPool(ctx context.Context) (*pgxpool.Pool, error) {
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

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

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