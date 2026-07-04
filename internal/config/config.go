// Package config carrega a configuração do serviço a partir de variáveis de ambiente.
package config

import (
	"fmt"
	"os"
)

// Config reúne todos os parâmetros de execução do serviço.
type Config struct {
	GRPCPort    string
	MetricsPort string
	DatabaseURL string
	LogLevel    string
}

// Load lê as variáveis de ambiente; DATABASE_URL tem prioridade sobre as PG*.
func Load() Config {
	c := Config{
		GRPCPort:    getenv("GRPC_PORT", "50051"),
		MetricsPort: getenv("METRICS_PORT", "9090"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    getenv("LOG_LEVEL", "info"),
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
