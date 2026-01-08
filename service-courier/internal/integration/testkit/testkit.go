package testkit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestPostgres struct {
	Container *postgres.PostgresContainer
	DB        *pgxpool.Pool
	DSN       string
}

func StartPostgres(ctx context.Context) (*TestPostgres, error) {
	dbName := "testdb"
	dbUser := "test"
	dbPassword := "test"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := mustInitPool(ctx, dsn)
    if err != nil {
        log.Fatal("error connecting to database: ", err)
    }

	return &TestPostgres{
		Container: postgresContainer,
		DB: db,
		DSN: dsn,
	}, nil
}

func (tp *TestPostgres) Close(ctx context.Context) {
	tp.DB.Close()
	_ = tp.Container.Terminate(ctx)
}

func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}


func mustInitPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
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