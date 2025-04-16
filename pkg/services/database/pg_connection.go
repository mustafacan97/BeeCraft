package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitializePgxConnectionPool() *pgxpool.Pool {
	// Load connection string from environment variable.
	connStr := os.Getenv("DB_CONNECTION_STRING")

	// Parse the connection string into a pgxpool configuration.
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(fmt.Errorf("database connection config parsing failed: %w", err))
	}

	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 15
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Second * 5

	config.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		zap.L().Info("Before acquiring the connection pool to the database")
		return true
	}

	config.AfterRelease = func(c *pgx.Conn) bool {
		zap.L().Info("After releasing the connection pool to the database")
		return true
	}

	config.BeforeClose = func(c *pgx.Conn) {
		zap.L().Info("Closed the connection pool to the database")
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Errorf("database connection failed: %v", err))
	}
	zap.L().Info("Successfully connected to PostgreSQL database.")
	return connPool
}
