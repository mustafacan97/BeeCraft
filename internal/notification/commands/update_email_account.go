package commands

import (
	"context"
	"errors"
	internalDomain "platform/internal/notification/domain"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/pkg/domain/valueobject"

	"github.com/google/uuid"
)

type UpdateEmailAccountCommand struct {
	ID           uuid.UUID
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

type UpdateEmailAccountCommandResponse struct{}

type UpdateEmailAccountCommandHandler struct {
	repository repositories.EmailAccountRepository
}

func NewUpdateEmailAccountCommandHandler(repository repositories.EmailAccountRepository) *UpdateEmailAccountCommandHandler {
	return &UpdateEmailAccountCommandHandler{repository: repository}
}

func (c *UpdateEmailAccountCommandHandler) Handle(ctx context.Context, command *UpdateEmailAccountCommand) (*UpdateEmailAccountCommandResponse, error) {
	email, err := valueobject.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	ea, _ := c.repository.GetByID(ctx, command.ID)

	if !ea.GetEmail().Equals(email) {
		ea.SetEmail(email)
	}

	if ea.GetDisplayName() != command.DisplayName {
		ea.SetDisplayName(command.DisplayName)
	}

	if ea.GetHost() != command.Host {
		ea.SetHost(command.Host)
	}

	if ea.GetPort() != command.Port {
		ea.SetPort(command.Port)
	}

	if ea.GetEnableSSL() != command.EnableSSL {
		ea.SetEnableSSL(command.EnableSSL)
	}

	if ea.GetSmtpType() != command.TypeID {
		ea.SetSMTPType(command.TypeID)
	}

	switch command.TypeID {
	case internalDomain.Login:
		encrypted, err := internalValueObject.EncryptAES(command.Password)
		if err != nil {
			return nil, err
		}
		credentials, err := internalValueObject.NewTraditionalCredentials(command.Username, encrypted)
		if err != nil {
			return nil, err
		}
		if ea.TraditionalCredentials == nil || !ea.TraditionalCredentials.Equals(credentials) {
			ea.SetTraditionalCredentials(credentials)
		}
	case internalDomain.GmailOAuth2:
	case internalDomain.MicrosoftOAuth2:
		credentials, err := internalValueObject.NewOAuth2Credentials(command.ClientID, command.ClientSecret, command.TenantID)
		if err != nil {
			return nil, err
		}
		if ea.OAuth2Credentials == nil || !ea.OAuth2Credentials.Equals(credentials) {
			ea.SetOAuth2Credentials(credentials)
		}
	default:
		return nil, errors.New("invalid login type")
	}

	err = c.repository.Update(ctx, ea)
	if err != nil {
		return nil, err
	}

	return &UpdateEmailAccountCommandResponse{}, nil
}
