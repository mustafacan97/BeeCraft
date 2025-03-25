package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "platform/config"
	"platform/internal/application/adapters/secondary/postgresql"
	"platform/internal/application/handlers"
	"platform/internal/application/handlers/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	zap.L().Info("Application is starting...")

	// Ensure logs are flushed before the application exits
	defer func() {
		// The ENOTTY error occurs when an operation that requires a terminal device is performed
		// on a file descriptor that is not connected to a terminal. We can ignore it for now.
		if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			zap.L().Error("Failed to flush logs", zap.Error(err))
		}
	}()

	// Initialize PostgreSQL connection pool
	dbPool := createPgxConnectionPool()
	defer dbPool.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	})

	// Initialize default config
	app.Use(pprof.New())

	// PostgreSQL Repositories
	userRepository := postgresql.NewUserRepository(dbPool)
	roleRepository := postgresql.NewRoleRepository(dbPool)

	// Handlers
	registerHandler := auth.NewRegisterHandler(userRepository, roleRepository)

	// Routes
	app.Post("/register", handlers.Serve(registerHandler))

	// Define a sample route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Run server in a goroutine
	go func() {
		zap.L().Info("Server is running on port 3000")
		if err := app.Listen(":3000"); err != nil {
			zap.L().Error("An error occurred on application starting %v", zap.Error(err))
		}
	}()

	// Graceful shutdown handling
	gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
	// Create a buffered channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Block until signal is received
	zap.L().Info("Shutting down server")

	// Create a context with timeout for graceful shutdown (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown Fiber server
	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.L().Error("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exited gracefully")
}

func createPgxConnectionPool() *pgxpool.Pool {
	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig())
	if err != nil {
		panic(fmt.Errorf("database connection failed: %v", err))
	}

	return connPool
}

func pgxPoolConfig() *pgxpool.Config {
	const DATABASE_URL string = "postgres://admin:123qwe@localhost:6432/platform?sslmode=disable"

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
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
