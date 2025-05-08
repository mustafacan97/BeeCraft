package commands

import (
	"context"
	"fmt"
	"platform/internal/notification/repositories"
	email_sender "platform/internal/notification/services/emailSender"
	"platform/internal/notification/services/encryption"
	voExternal "platform/pkg/domain/value_object"
)

type SendTestEmailCommand struct {
	From string
	To   string
}

type SendTestEmailCommandResponse struct{}

type SendTestEmailCommandHandler struct {
	encryption encryption.EncryptionService
	repository repositories.EmailAccountRepository
}

func NewSendTestEmailCommandHandler(encryption encryption.EncryptionService, repository repositories.EmailAccountRepository) *SendTestEmailCommandHandler {
	return &SendTestEmailCommandHandler{
		encryption: encryption,
		repository: repository,
	}
}

func (c *SendTestEmailCommandHandler) Handle(ctx context.Context, command *SendTestEmailCommand) (*SendTestEmailCommandResponse, error) {
	fromEmail, err := voExternal.NewEmail(command.From)
	if err != nil {
		return nil, err
	}

	toEmail, err := voExternal.NewEmail(command.To)
	if err != nil {
		return nil, err
	}

	ea, _ := c.repository.GetByEmail(ctx, fromEmail)
	if ea == nil {
		return nil, fmt.Errorf("email not found: %s", command.From)
	}

	email, _ := email_sender.BaseEmailDetail("Test Email", "<h1>Hello World!</h1>", fromEmail, toEmail)
	email_sender.SendEmail(c.encryption, ea, email)
	return nil, nil
}
