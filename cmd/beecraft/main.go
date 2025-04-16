package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"platform/pkg/services/database"
	"platform/pkg/services/eventbus"
	"platform/pkg/services/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logging.InitializeLogging()

	zap.L().Info("Application is starting...")

	// Ensure logs are flushed before the application exits
	defer func() {
		// The ENOTTY error occurs when an operation that requires a terminal device is performed
		// on a file descriptor that is not connected to a terminal. We can ignore it for now.
		if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			zap.L().Error("Failed to flush logs", zap.Error(err))
		}
	}()

	if err := godotenv.Load(); err != nil {
		zap.L().Fatal("Error loading .env file", zap.Error(err))
	}

	// Initialize rabbitMQ configurations
	bus := eventbus.NewRabbitMQEventBus("amq.topic")
	defer bus.Close()

	// Initialize database configurations
	dbPool := database.InitializePgxConnectionPool()
	defer dbPool.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		ErrorHandler: errorHandler(zap.L()),
	})

	// Middlewares
	app.Use(zapLoggerMiddleware(zap.L()))
	app.Use(recover.New())
	app.Use(pprof.New()) // Enable pprof middleware for performance profiling and debugging

	SetupRouter(app, dbPool, bus)

	go func() {
		zap.L().Info("Server is running on port 3000")
		if err := app.Listen(":3000"); err != nil {
			zap.L().Error("An error occurred on application starting", zap.Error(err))
		}
	}()

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

func errorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		logger.Error("request error",
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Error(err),
		)

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}

func zapLoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		// Structured log
		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", stop.Sub(start)),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}
