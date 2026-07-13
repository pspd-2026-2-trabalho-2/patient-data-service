// Package config carrega a configuração do serviço a partir de variáveis de ambiente.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config reúne todos os parâmetros de execução do serviço.
type Config struct {
	GRPCPort    string
	MetricsPort string
	DatabaseURL string
	LogLevel    string
	DBMaxConns  int32
	DBMinConns  int32
	RPCTimeout  time.Duration
}

// Load lê as variáveis de ambiente; DATABASE_URL tem prioridade sobre as PG*.
func Load() Config {
	c := Config{
		GRPCPort:    getenv("GRPC_PORT", "50051"),
		MetricsPort: getenv("METRICS_PORT", "9090"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    getenv("LOG_LEVEL", "info"),
		DBMaxConns:  int32(getenvInt("DB_MAX_CONNS", 10)),
		DBMinConns:  int32(getenvInt("DB_MIN_CONNS", 1)),
		RPCTimeout:  getenvDuration("RPC_TIMEOUT", 30*time.Second),
	}
	if c.DatabaseURL == "" {
		c.DatabaseURL = buildDSN()
	}
	return c
}

func buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getenv("PGUSER", "pspd"),
		getenv("PGPASSWORD", "pspd"),
		getenv("PGHOST", "localhost"),
		getenv("PGPORT", "5432"),
		getenv("PGDATABASE", "hospital"),
		getenv("PGSSLMODE", "disable"),
	)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
