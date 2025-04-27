package commands

import (
	"context"
	"errors"
	internalDomain "platform/internal/notification/domain"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/internal/shared"
	"platform/pkg/domain/valueobject"

	"github.com/google/uuid"
)

type CreateEmailAccountCommand struct {
	Email        string
	DisplayName  string
	Host         string
	Port         int
	Username     string
	Password     string
	EnableSSL    bool
	TypeID       int
	ClientID     string
	TenantID     string
	ClientSecret string
}

type CreateEmailAccountCommandResponse struct{}

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

	projectIDContext := ctx.Value(shared.ProjectIDContextKey).(string)
	projectID, err := uuid.Parse(projectIDContext)
	if err != nil {
		return nil, err
	}

	emailAccount := internalDomain.NewEmailAccount(uuid.New(), projectID, command.TypeID, email, command.DisplayName, command.Host, command.Port, command.EnableSSL)

	if command.TypeID == internalDomain.Login {
		encrypted, err := internalValueObject.EncryptPassword(command.Password, []byte(command.Email))
		if err == nil {
			emailAccount.TraditionalCredentials = internalValueObject.NewTraditionalCredentials(command.Username, *encrypted)
		} else {
			return nil, err
		}
	} else if command.TypeID == internalDomain.GmailOAuth2 || command.TypeID == internalDomain.MicrosoftOAuth2 {
		emailAccount.OAuth2Credentials = internalValueObject.NewOAuth2Credentials(command.ClientID, command.ClientSecret, command.TenantID)
	} else {
		return nil, errors.New("unsupported SMTP type")
	}

	err = c.repository.Create(ctx, emailAccount)
	if err != nil {
		return nil, err
	}

	return &CreateEmailAccountCommandResponse{}, nil
}
