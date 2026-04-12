package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                     string
	ServerPort                 int
	DBConnectionString         string
	DeliveryCheckerInterval    time.Duration
	EnablePprof                bool
	PprofAddr                  string
	DBMaxConns                 int32
	DBMinConns                 int32
	DBMaxConnLifetime          time.Duration
}

func Load() (Config, error) {
	_ = loadDotEnvIfNeeded()

	cfg := Config{
		AppEnv:                  getString("APP_ENV", "local"),
		ServerPort:              getInt("SERVER_PORT", 8080),
		DBConnectionString:      strings.TrimSpace(os.Getenv("DB_CONNECTION_STRING")),
		DeliveryCheckerInterval: time.Duration(getInt("DELIVERY_CHECKER_INTERVAL", 10)) * time.Second,
		EnablePprof:             getBool("PPROF_ENABLED", false),
		PprofAddr:               getString("PPROF_ADDR", "localhost:9000"),
		DBMaxConns:              int32(getInt("DB_MAX_CONNS", 10)),
		DBMinConns:              int32(getInt("DB_MIN_CONNS", 5)),
		DBMaxConnLifetime:       time.Duration(getInt("DB_MAX_CONN_LIFETIME_SEC", 3600)) * time.Second,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.DBConnectionString == "" {
		return fmt.Errorf("DB_CONNECTION_STRING is required")
	}
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("SERVER_PORT must be in range 1..65535")
	}
	if c.DeliveryCheckerInterval <= 0 {
		return fmt.Errorf("DELIVERY_CHECKER_INTERVAL must be > 0")
	}
	if c.DBMinConns < 0 || c.DBMaxConns <= 0 || c.DBMinConns > c.DBMaxConns {
		return fmt.Errorf("invalid db pool config: min=%d max=%d", c.DBMinConns, c.DBMaxConns)
	}
	return nil
}

func loadDotEnvIfNeeded() error {
	if os.Getenv("APP_ENV") == "docker" {
		return nil
	}
	return godotenv.Load()
}

func getString(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
