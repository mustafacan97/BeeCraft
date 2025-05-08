package queries

import (
	"context"
	voInternal "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	voExternal "platform/pkg/domain/value_object"
	"time"

	"github.com/google/uuid"
)

type GetEmailAccountByEmailQuery struct {
	Email string
}

type GetEmailAccountByEmailQueryResponse struct {
	ID                     uuid.UUID
	ProjectID              uuid.UUID
	Email                  voExternal.Email
	DisplayName            string
	Host                   string
	Port                   int
	EnableSSL              bool
	TypeId                 int
	TraditionalCredentials *voInternal.TraditionalCredentials
	OAuth2Credentials      *voInternal.OAuth2Credentials
	CreatedAt              time.Time
}

type GetEmailAccountByEmailQueryHandler struct {
	repository repositories.EmailAccountRepository
}

func NewGetEmailAccountByEmailQueryHandler(repository repositories.EmailAccountRepository) *GetEmailAccountByEmailQueryHandler {
	return &GetEmailAccountByEmailQueryHandler{repository: repository}
}

func (c *GetEmailAccountByEmailQueryHandler) Handle(ctx context.Context, query *GetEmailAccountByEmailQuery) (*GetEmailAccountByEmailQueryResponse, error) {
	email, err := voExternal.NewEmail(query.Email)
	if err != nil {
		return nil, err
	}

	emailAccount, err := c.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if emailAccount == nil {
		return nil, nil
	}

	return &GetEmailAccountByEmailQueryResponse{
		ID:                     emailAccount.GetID(),
		ProjectID:              emailAccount.GetProjectID(),
		Email:                  emailAccount.GetEmail(),
		DisplayName:            emailAccount.GetDisplayName(),
		Host:                   emailAccount.GetHost(),
		Port:                   emailAccount.GetPort(),
		EnableSSL:              emailAccount.GetEnableSSL(),
		TypeId:                 emailAccount.GetSmtpType(),
		TraditionalCredentials: emailAccount.GetTraditionalCredentials(),
		OAuth2Credentials:      emailAccount.GetOAuth2Credentials(),
		CreatedAt:              emailAccount.GetCreatedAt(),
	}, nil
}
