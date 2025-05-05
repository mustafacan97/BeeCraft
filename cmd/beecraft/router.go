package main

import (
	baseHandler "platform/internal/shared/handlers"
	"platform/internal/shared/middlewares"
	event_bus "platform/pkg/services/eventbus"
	mediator "platform/pkg/services/mediator"

	iamHandlers "platform/internal/iam/handlers"
	"platform/internal/notification/commands"
	notificationHandlers "platform/internal/notification/handlers"
	"platform/internal/notification/queries"

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
		createHandler := &notificationHandlers.CreateEmailAccountHandler{}
		updateHandler := &notificationHandlers.UpdateEmailAccountHandler{}
		deleteHandler := &notificationHandlers.DeleteEmailAccountHandler{}
		getHandler := &notificationHandlers.GetEmailAccountHandler{}
		listHandler := &notificationHandlers.ListEmailAccountHandler{}

		notificationGroup.Post("/email-account", baseHandler.Serve(createHandler))
		notificationGroup.Put("/email-account/:id", baseHandler.Serve(updateHandler))
		notificationGroup.Delete("/email-account/:id", baseHandler.Serve(deleteHandler))
		notificationGroup.Get("/email-account/:id", baseHandler.Serve(getHandler))
		notificationGroup.Get("/email-accounts", baseHandler.Serve(listHandler))

		oauthUrlHandler := notificationHandlers.NewOAuthUrlHandler()
		oauthCallbackHandler := notificationHandlers.NewOAuthCallbackHandler(&emailAccountRepository)
		notificationGroup.Post("/oauth", baseHandler.Serve(oauthUrlHandler))
		notificationGroup.Get("/oauth-callback", baseHandler.Serve(oauthCallbackHandler))
	}
}
