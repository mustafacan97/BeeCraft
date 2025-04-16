package main

import (
	baseHandler "platform/internal/shared/handlers"
	event_bus "platform/pkg/services/eventbus"

	iamHandlers "platform/internal/iam/handlers"
	notificationHandlers "platform/internal/notification/handlers"

	iamRepositories "platform/internal/iam/repositories"
	notificationRepositories "platform/internal/notification/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRouter configures the Fiber app with Zap logging, recovery, routes, and handlers.
func SetupRouter(app *fiber.App, dbPool *pgxpool.Pool, bus event_bus.EventBus) {
	// PostgreSQL Repositories
	userRepository := iamRepositories.NewUserRepository(dbPool)
	roleRepository := iamRepositories.NewRoleRepository(dbPool)
	emailAccountRepository := notificationRepositories.NewPgEmailAccountRepository(dbPool)

	// Handlers
	registerHandler := iamHandlers.NewRegisterHandler(bus, &userRepository, &roleRepository)
	oauthUrlHandler := notificationHandlers.NewOAuthUrlHandler()
	oauthCallbackHandler := notificationHandlers.NewOAuthCallbackHandler(&emailAccountRepository)

	// Health-check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// API Versioning
	version1 := app.Group("/v1")

	// IAM Service Routes
	iamGroup := version1.Group("/iam")
	{
		iamGroup.Post("/register", baseHandler.Serve(registerHandler))
	}

	// Notification Service Routes
	notificationGroup := version1.Group("/notification")
	{
		notificationGroup.Post("/oauth", baseHandler.Serve(oauthUrlHandler))
		notificationGroup.Get("/oauth-callback", baseHandler.Serve(oauthCallbackHandler))
	}
}
