package main

import (
	baseHandler "platform/internal/shared/handlers"
	"platform/internal/shared/middlewares"
	"platform/pkg/services/cache"
	event_bus "platform/pkg/services/eventbus"
	mediator "platform/pkg/services/mediator"

	iamHandlers "platform/internal/iam/handlers"
	notificationHandlers "platform/internal/notification/handlers"
	"platform/internal/notification/mediatr/commands"
	event_notification "platform/internal/notification/mediatr/notifications"
	"platform/internal/notification/mediatr/queries"

	iamRepositories "platform/internal/iam/repositories"
	notificationRepositories "platform/internal/notification/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRouter configures the Fiber app with Zap logging, recovery, routes, and handlers.
func SetupRouter(app *fiber.App, dbPool *pgxpool.Pool, bus event_bus.EventBus) {
	// Cache service
	cacheService := cache.NewMemcacheManager("localhost:11211")

	// PostgreSQL Repositories
	userRepository := iamRepositories.NewUserRepository(dbPool)
	roleRepository := iamRepositories.NewRoleRepository(dbPool)
	emailAccountRepository := notificationRepositories.NewPgEmailAccountRepository(dbPool, cacheService)

	// Mediator
	getEmailAccountByEmailQueryHandler := queries.NewGetEmailAccountByEmailQueryHandler(emailAccountRepository)
	getEmailAccountByIDQueryHandler := queries.NewGetEmailAccountByIDQueryHandler(emailAccountRepository)
	listEmailAccountQueryHandler := queries.NewListEmailAccountQueryHandler(emailAccountRepository)
	createEmailAccountCommandHandler := commands.NewCreateEmailAccountCommandHandler(emailAccountRepository)
	updateEmailAccountCommandHandler := commands.NewUpdateEmailAccountCommandHandler(emailAccountRepository)
	deleteEmailAccountCommandHandler := commands.NewDeleteEmailAccountCommandHandler(emailAccountRepository)
	mediator.RegisterRequestHandler(getEmailAccountByEmailQueryHandler)
	mediator.RegisterRequestHandler(getEmailAccountByIDQueryHandler)
	mediator.RegisterRequestHandler(listEmailAccountQueryHandler)
	mediator.RegisterRequestHandler(createEmailAccountCommandHandler)
	mediator.RegisterRequestHandler(updateEmailAccountCommandHandler)
	mediator.RegisterRequestHandler(&deleteEmailAccountCommandHandler)

	// Notification Handlers
	emailAccountCreatedHandler := event_notification.EmailAccountCreatedEventHandler{}
	mediator.RegisterNotificationHandler(&emailAccountCreatedHandler)

	// Handlers
	registerHandler := iamHandlers.NewRegisterHandler(bus, &userRepository, &roleRepository)

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
	notificationGroup := version1.Group("/notification", middlewares.RequireProjectID())
	{
		createHandler := notificationHandlers.CreateEmailAccountHandler{}
		updateHandler := &notificationHandlers.UpdateEmailAccountHandler{}
		deleteHandler := &notificationHandlers.DeleteEmailAccountHandler{}
		getHandler := &notificationHandlers.GetEmailAccountHandler{}
		listHandler := &notificationHandlers.ListEmailAccountHandler{}
		oauth2CallbackHandler := &notificationHandlers.OAuth2CallbackHandler{}
		testEmailHandler := &notificationHandlers.SendTestEmailHandler{}

		notificationGroup.Post("/email-accounts", baseHandler.Serve(&createHandler))
		notificationGroup.Put("/email-account/:id", baseHandler.Serve(updateHandler))
		notificationGroup.Delete("/email-account/:id", baseHandler.Serve(deleteHandler))
		notificationGroup.Get("/email-account/:id", baseHandler.Serve(getHandler))
		notificationGroup.Get("/email-accounts", baseHandler.Serve(listHandler))
		notificationGroup.Get("/oauth2-callback", baseHandler.Serve(oauth2CallbackHandler))
		notificationGroup.Get("/email-accounts/:id/:email", baseHandler.Serve(testEmailHandler))
	}
}
