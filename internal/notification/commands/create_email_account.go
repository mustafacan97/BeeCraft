package commands

import (
	"context"
	internalDomain "platform/internal/notification/domain"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/internal/shared"
	"platform/pkg/domain/valueobject"

	"github.com/google/uuid"
)

type CreateEmailAccountCommand struct {
	ProjectID    uuid.UUID
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
	ID uuid.UUID
}

type CreateEmailAccountCommandHandler struct {
	repository repositories.EmailAccountRepository
}

func NewCreateEmailAccountCommandHandler(repository repositories.EmailAccountRepository) *CreateEmailAccountCommandHandler {
	return &CreateEmailAccountCommandHandler{repository: repository}
}

func (c *CreateEmailAccountCommandHandler) Handle(ctx context.Context, command *CreateEmailAccountCommand) (*CreateEmailAccountCommandResponse, error) {
	email, err := valueobject.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	projectID := ctx.Value(shared.ProjectIDContextKey).(uuid.UUID)
	emailAccount := internalDomain.NewEmailAccount(uuid.New(), projectID, command.TypeID, email, command.DisplayName, command.Host, command.Port, command.EnableSSL)

	if command.TypeID == internalDomain.Login {
		encrypted, err := internalValueObject.EncryptAES(command.Password)
		if err != nil {
			return nil, err
		}
		credentials, err := internalValueObject.NewTraditionalCredentials(command.Username, encrypted)
		if err != nil {
			return nil, err
		}
		emailAccount.SetTraditionalCredentials(credentials)
	} else {
		credentials, err := internalValueObject.NewOAuth2Credentials(command.ClientID, command.ClientSecret, command.TenantID)
		if err != nil {
			return nil, err
		}
		emailAccount.SetOAuth2Credentials(credentials)
	}

	err = c.repository.Create(ctx, emailAccount)
	if err != nil {
		return nil, err
	}

	return &CreateEmailAccountCommandResponse{ID: emailAccount.GetID()}, nil
}
