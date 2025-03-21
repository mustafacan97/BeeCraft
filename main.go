package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "platform/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	zap.L().Info("Application is starting...")

	// Ensure logs are flushed before the application exits
	defer func() {
		if err := zap.L().Sync(); err != nil {
			zap.L().Error("Failed to flush logs", zap.Error(err))
		}
	}()

	dbConn, err := connectToPostgresql()
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}
	defer dbConn.Close(context.Background())

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	})

	// Initialize default config
	app.Use(pprof.New())

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

func connectToPostgresql() (*pgx.Conn, error) {
	dsn := "postgres://admin:123qwe@localhost:6432/platform?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	zap.L().Info("Connected to PostgreSQL")
	return conn, nil
}
