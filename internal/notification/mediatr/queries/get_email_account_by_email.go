package queries

import (
	"context"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/pkg/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

type GetEmailAccountByEmailQuery struct {
	Email string
}

type GetEmailAccountByEmailQueryResponse struct {
	ID                     uuid.UUID
	ProjectID              uuid.UUID
	Email                  valueobject.Email
	DisplayName            string
	Host                   string
	Port                   int
	EnableSSL              bool
	TypeId                 int
	TraditionalCredentials *internalValueObject.TraditionalCredential
	OAuth2Credentials      *internalValueObject.OAuth2Credential
	CreatedAt              time.Time
}

type GetEmailAccountByEmailQueryHandler struct {
	repository repositories.EmailAccountRepository
}

func NewGetEmailAccountByEmailQueryHandler(repository repositories.EmailAccountRepository) *GetEmailAccountByEmailQueryHandler {
	return &GetEmailAccountByEmailQueryHandler{repository: repository}
}

func (c *GetEmailAccountByEmailQueryHandler) Handle(ctx context.Context, query *GetEmailAccountByEmailQuery) (*GetEmailAccountByEmailQueryResponse, error) {
	email, err := valueobject.NewEmail(query.Email)
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

	base := &GetEmailAccountByEmailQueryResponse{
		ID:                     emailAccount.GetID(),
		ProjectID:              emailAccount.GetProjectID(),
		Email:                  emailAccount.GetEmail(),
		DisplayName:            emailAccount.GetDisplayName(),
		Host:                   emailAccount.GetHost(),
		Port:                   emailAccount.GetPort(),
		EnableSSL:              emailAccount.IsSslEnabled(),
		TypeId:                 emailAccount.GetSmtpType(),
		TraditionalCredentials: emailAccount.TraditionalCredentials,
		OAuth2Credentials:      emailAccount.OAuth2Credentials,
		CreatedAt:              emailAccount.GetCreatedDate(),
	}

	return base, nil
}
