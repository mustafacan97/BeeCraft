package commands

import (
	"context"
	"platform/internal/notification/domain"
	voInternal "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/internal/notification/services/encryption"
	voExternal "platform/pkg/domain/value_object"
	"time"
)

type UpdateEmailAccountCommand struct {
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
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
}

type UpdateEmailAccountCommandResponse struct{}

type UpdateEmailAccountCommandHandler struct {
	encryption encryption.EncryptionService
	repository repositories.EmailAccountRepository
}

func NewUpdateEmailAccountCommandHandler(encryption encryption.EncryptionService, repository repositories.EmailAccountRepository) *UpdateEmailAccountCommandHandler {
	return &UpdateEmailAccountCommandHandler{
		encryption: encryption,
		repository: repository,
	}
}

func (c *UpdateEmailAccountCommandHandler) Handle(ctx context.Context, command *UpdateEmailAccountCommand) (*UpdateEmailAccountCommandResponse, error) {
	email, err := voExternal.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	ea, _ := c.repository.GetByEmail(ctx, email)
	if ea == nil {
		return nil, nil
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
		ea.SetSmtpType(command.TypeID)
	}

	if command.TypeID == domain.Login {
		encrypted, err := c.encryption.Encrypt(command.Password)
		if err != nil {
			return nil, err
		}
		oldCredentials := ea.GetTraditionalCredentials()
		newCredentials := voInternal.NewTraditionalCredentials(command.Username, encrypted)
		if (oldCredentials == nil && newCredentials != nil) || !oldCredentials.Equals(newCredentials) {
			ea.SetTraditionalCredentials(newCredentials)
		}
	} else if command.TypeID == domain.GmailOAuth2 || command.TypeID == domain.MicrosoftOAuth2 {
		oldCredentials := ea.GetOAuth2Credentials()
		newCredentials := voInternal.NewOAuth2Credentials(command.ClientID, command.TenantID, command.ClientSecret)
		if (oldCredentials == nil && newCredentials != nil) || !oldCredentials.Equals(newCredentials) {
			ea.SetOAuth2Credentials(newCredentials)
		}

		oldTokenInfo := ea.GetTokenInformation()
		newTokenInfo := voInternal.NewTokenInformation(command.AccessToken, command.RefreshToken, command.ExpireAt)
		if (oldTokenInfo == nil && newTokenInfo != nil) || !oldTokenInfo.Equals(newTokenInfo) {
			ea.SetTokenInformation(newTokenInfo)
		}
	}

	err = c.repository.Update(ctx, ea)
	if err != nil {
		return nil, err
	}

	return &UpdateEmailAccountCommandResponse{}, nil
}
