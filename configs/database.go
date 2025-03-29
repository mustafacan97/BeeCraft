package configs

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
	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig())
	if err != nil {
		panic(fmt.Errorf("database connection failed: %v", err))
	}
	zap.L().Info("Successfully connected to PostgreSQL database.")
	return connPool
}

func pgxPoolConfig() *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		panic(fmt.Errorf("database connection config parsing failed: %w", err))
	}

	dbConfig.MaxConns = 100
	dbConfig.MinConns = 10
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = time.Minute * 15
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		zap.L().Info("Before acquiring the connection pool to the database")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		zap.L().Info("After releasing the connection pool to the database")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		zap.L().Info("Closed the connection pool to the database")
	}

	return dbConfig
}
