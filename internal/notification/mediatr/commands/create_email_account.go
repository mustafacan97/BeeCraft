package commands

import (
	"context"
	"platform/internal/notification/domain"
	voInternal "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/internal/notification/services/encryption"
	"platform/internal/shared"
	voExternal "platform/pkg/domain/value_object"

	"github.com/google/uuid"
)

type CreateEmailAccountCommand struct {
	Email        string
	DisplayName  string
	Host         string
	Port         int
	EnableSSL    bool
	TypeID       int
	Username     string
	Password     string
	ClientID     string
	TenantID     string
	ClientSecret string
}

type CreateEmailAccountCommandResponse struct {
}

type CreateEmailAccountCommandHandler struct {
	encryption encryption.EncryptionService
	repository repositories.EmailAccountRepository
}

func NewCreateEmailAccountCommandHandler(encryption encryption.EncryptionService, repository repositories.EmailAccountRepository) *CreateEmailAccountCommandHandler {
	return &CreateEmailAccountCommandHandler{
		encryption: encryption,
		repository: repository,
	}
}

func (c *CreateEmailAccountCommandHandler) Handle(ctx context.Context, command *CreateEmailAccountCommand) (*CreateEmailAccountCommandResponse, error) {
	// STEP-1: Get project identifier and validate
	pidVal := ctx.Value(shared.ProjectIDContextKey)
	projectID, ok := pidVal.(uuid.UUID)
	if !ok {
		return nil, shared.ErrInvalidContext
	}

	email, err := voExternal.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	ea := domain.EmailAccount{}
	emailAccountID := uuid.New()
	ea.SetID(emailAccountID)
	ea.SetProjectID(projectID)
	ea.SetEmail(email)
	ea.SetDisplayName(command.DisplayName)
	ea.SetHost(command.Host)
	ea.SetPort(command.Port)
	ea.SetEnableSSL(command.EnableSSL)
	ea.SetSmtpType(command.TypeID)

	if command.TypeID == domain.Login {
		encrypted, err := c.encryption.Encrypt(command.Password)
		if err != nil {
			return nil, err
		}
		credentials := voInternal.NewTraditionalCredentials(command.Username, encrypted)
		ea.SetTraditionalCredentials(credentials)
	} else {
		credentials := voInternal.NewOAuth2Credentials(command.ClientID, command.TenantID, command.ClientSecret)
		ea.SetOAuth2Credentials(credentials)
	}

	err = c.repository.Create(ctx, &ea)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
