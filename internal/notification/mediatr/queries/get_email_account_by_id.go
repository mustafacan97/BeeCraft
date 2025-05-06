package queries

import (
	"context"
	voInternal "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	voExternal "platform/pkg/domain/value_object"
	"time"

	"github.com/google/uuid"
)

type GetEmailAccountByIDQuery struct {
	ID uuid.UUID
}

type GetEmailAccountByIDQueryResponse struct {
	ID                     uuid.UUID                          `json:"id"`
	ProjectID              uuid.UUID                          `json:"project_id"`
	TypeId                 int                                `json:"type_id"`
	DisplayName            string                             `json:"display_name"`
	Host                   string                             `json:"host"`
	EnableSSL              bool                               `json:"enable_ssl"`
	CreatedAt              time.Time                          `json:"created_at"`
	Port                   int                                `json:"port"`
	Email                  voExternal.Email                   `json:"email"`
	TraditionalCredentials *voInternal.TraditionalCredentials `json:"traditional_credentials"`
	OAuth2Credentials      *voInternal.OAuth2Credentials      `json:"oauth2_credentials"`
	TokenInformation       *voInternal.TokenInformation       `json:"token_information"`
}

type GetEmailAccountByIDQueryHandler struct {
	repository repositories.EmailAccountRepository
}

func NewGetEmailAccountByIDQueryHandler(repository repositories.EmailAccountRepository) *GetEmailAccountByIDQueryHandler {
	return &GetEmailAccountByIDQueryHandler{repository: repository}
}

func (c *GetEmailAccountByIDQueryHandler) Handle(ctx context.Context, query *GetEmailAccountByIDQuery) (*GetEmailAccountByIDQueryResponse, error) {
	emailAccount, err := c.repository.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	if emailAccount == nil {
		return nil, nil
	}

	base := &GetEmailAccountByIDQueryResponse{
		ID:                     emailAccount.ID,
		ProjectID:              emailAccount.GetProjectID(),
		Email:                  emailAccount.GetEmail(),
		DisplayName:            emailAccount.GetDisplayName(),
		Host:                   emailAccount.GetHost(),
		Port:                   emailAccount.GetPort(),
		EnableSSL:              emailAccount.GetEnableSSL(),
		TypeId:                 emailAccount.GetSmtpType(),
		TraditionalCredentials: emailAccount.GetTraditionalCredentials(),
		OAuth2Credentials:      emailAccount.GetOAuth2Credentials(),
		TokenInformation:       emailAccount.GetTokenInformation(),
		CreatedAt:              emailAccount.GetCreatedAt(),
	}

	return base, nil
}
