package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"platform/configs"
	"platform/internal/application/adapters/postgresql"
	"platform/internal/application/handlers"
	baseHandler "platform/internal/shared/handlers"
	eventBusAdapter "platform/pkg/services/eventbus"

	notificationHandlers "platform/internal/notification/handlers"

	notificationRepositories "platform/internal/notification/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	configs.InitializeLogConfig()

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
		zap.L().Fatal("Error loading .env file")
	}
	zap.L().Info(".env file loaded successfully")

	// For InMemory
	// bus := eventBusAdapter.NewInMemoryEventBus()
	bus := eventBusAdapter.NewRabbitMQEventBus("amq.topic")
	defer bus.Close()
	/*bus.Subscribe("NotificationService", "user.registered", func(ctx context.Context, event domain.Event) error {
		zap.L().Info("Received event", zap.Any("event", event))
		return nil
	})*/

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	})

	// Initialize default config
	app.Use(pprof.New())

	// Initialize database configurations
	dbPool := configs.InitializePgxConnectionPool()
	defer dbPool.Close()

	// PostgreSQL Repositories
	userRepository := postgresql.NewUserRepository(dbPool)
	roleRepository := postgresql.NewRoleRepository(dbPool)
	emailAccountRepository := notificationRepositories.NewPgEmailAccountRepository(dbPool)

	// Handlers
	registerHandler := handlers.NewRegisterHandler(&bus, &userRepository, &roleRepository)
	oauthUrlHandler := notificationHandlers.NewOAuthUrlHandler()
	oauthCallbackHandler := notificationHandlers.NewOAuthCallbackHandler(&emailAccountRepository)

	// Routes
	app.Post("/register", handlers.Serve(registerHandler))
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("Hello, World!") })
	app.Post("/oauth", baseHandler.Serve(oauthUrlHandler))
	app.Get("/oauth-callback", baseHandler.Serve(oauthCallbackHandler))

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
