package commands

import (
	"context"
	"platform/internal/notification/repositories"
	email_sender "platform/internal/notification/services/emailSender"
	"platform/internal/notification/services/encryption"
	vo "platform/pkg/domain/value_object"

	"github.com/google/uuid"
)

type SendTestEmailCommand struct {
	EmailAccountID uuid.UUID
	To             string
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

func (c *SendTestEmailCommandHandler) Handle(ctx context.Context, command *SendTestEmailCommand) (*UpdateEmailAccountCommandResponse, error) {
	ea, _ := c.repository.GetByID(ctx, command.EmailAccountID)
	if ea == nil {
		return nil, nil
	}

	toEmail, _ := vo.NewEmail(command.To)
	email, _ := email_sender.BaseEmailDetail("Test Email", "<h1>Hello World!</h1>", ea.GetEmail(), toEmail)
	email_sender.SendEmail(c.encryption, ea, email)
	return nil, nil
}
