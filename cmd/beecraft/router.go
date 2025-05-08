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
	"platform/internal/notification/services/encryption"

	iamRepositories "platform/internal/iam/repositories"
	notificationRepositories "platform/internal/notification/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRouter configures the Fiber app with Zap logging, recovery, routes, and handlers.
func SetupRouter(app *fiber.App, dbPool *pgxpool.Pool, bus event_bus.EventBus) {
	// Services
	cacheService := cache.NewMemcacheManager("localhost:11211")
	encryptionService, _ := encryption.NewAESEncryptionService([]byte("1234567890123456"))

	// Repositories
	userRepository := iamRepositories.NewUserRepository(dbPool)
	roleRepository := iamRepositories.NewRoleRepository(dbPool)
	emailAccountRepository := notificationRepositories.NewPgEmailAccountRepository(dbPool, cacheService)

	// Mediator Queries
	getAllEmailAccountQueryHandler := queries.NewGetAllEmailAccountQueryHandler(emailAccountRepository)
	getEmailAccountByEmailQueryHandler := queries.NewGetEmailAccountByEmailQueryHandler(emailAccountRepository)
	mediator.RegisterRequestHandler(getAllEmailAccountQueryHandler)
	mediator.RegisterRequestHandler(getEmailAccountByEmailQueryHandler)

	// Mediator Commands
	createEmailAccountCommandHandler := commands.NewCreateEmailAccountCommandHandler(encryptionService, emailAccountRepository)
	deleteEmailAccountCommandHandler := commands.NewDeleteEmailAccountCommandHandler(emailAccountRepository)
	sendTestEmailCommandHandler := commands.NewSendTestEmailCommandHandler(encryptionService, emailAccountRepository)
	updateEmailAccountCommandHandler := commands.NewUpdateEmailAccountCommandHandler(encryptionService, emailAccountRepository)
	mediator.RegisterRequestHandler(createEmailAccountCommandHandler)
	mediator.RegisterRequestHandler(deleteEmailAccountCommandHandler)
	mediator.RegisterRequestHandler(sendTestEmailCommandHandler)
	mediator.RegisterRequestHandler(updateEmailAccountCommandHandler)

	// Notification Handlers
	emailAccountCreatedHandler := event_notification.EmailAccountCreatedEventHandler{}
	mediator.RegisterNotificationHandler(&emailAccountCreatedHandler)

	// Health-check
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// API Versioning
	version1 := app.Group("/v1")

	// IAM Service Routes
	iamGroup := version1.Group("/iam")
	{
		registerHandler := iamHandlers.NewRegisterHandler(bus, &userRepository, &roleRepository)
		iamGroup.Post("/register", baseHandler.Serve(registerHandler))
	}

	// Notification Service Routes
	notificationGroup := version1.Group("/notification", middlewares.ProjectIDInjector())
	{
		testEmailHandler := notificationHandlers.SendTestEmailHandler{}
		notificationGroup.Post("/email-accounts/:from", baseHandler.Serve(&testEmailHandler))

		createHandler := notificationHandlers.CreateEmailAccountHandler{}
		notificationGroup.Post("/email-accounts", baseHandler.Serve(&createHandler))

		deleteHandler := notificationHandlers.DeleteEmailAccountHandler{}
		notificationGroup.Delete("/email-accounts/:email", baseHandler.Serve(&deleteHandler))

		getAllHandler := notificationHandlers.GetAllEmailAccountHandler{}
		notificationGroup.Get("/email-accounts", baseHandler.Serve(&getAllHandler))

		oauth2CallbackHandler := notificationHandlers.OAuth2CallbackHandler{}
		notificationGroup.Get("/email-accounts/oauth2-callback", baseHandler.Serve(&oauth2CallbackHandler))

		getHandler := notificationHandlers.GetEmailAccountHandler{}
		notificationGroup.Get("/email-accounts/:email", baseHandler.Serve(&getHandler))

		updateHandler := notificationHandlers.UpdateEmailAccountHandler{}
		notificationGroup.Put("/email-accounts/:email", baseHandler.Serve(&updateHandler))
	}
}
