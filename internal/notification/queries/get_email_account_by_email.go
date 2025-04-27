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
	ID                     uuid.UUID                                  `json:"id"`
	ProjectID              uuid.UUID                                  `json:"project_id"`
	TypeId                 int                                        `json:"type_id"`
	DisplayName            string                                     `json:"display_name"`
	Host                   string                                     `json:"host"`
	EnableSSL              bool                                       `json:"enable_ssl"`
	CreatedAt              time.Time                                  `json:"created_at"`
	Port                   int                                        `json:"port"`
	Email                  valueobject.Email                          `json:"email"`
	TraditionalCredentials *internalValueObject.TraditionalCredential `json:"traditional_credentials"`
	OAuth2Credentials      *internalValueObject.OAuth2Credential      `json:"oauth2_credentials"`
	TokenInformation       *internalValueObject.TokenInformation      `json:"token_information"`
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
		ID:                     emailAccount.ID,
		ProjectID:              emailAccount.GetProjectID(),
		Email:                  emailAccount.Email,
		DisplayName:            emailAccount.GetDisplayName(),
		Host:                   emailAccount.GetHost(),
		Port:                   emailAccount.GetPort(),
		TraditionalCredentials: emailAccount.TraditionalCredentials,
		EnableSSL:              emailAccount.IsSslEnabled(),
		TypeId:                 emailAccount.GetSmtpType(),
		OAuth2Credentials:      emailAccount.OAuth2Credentials,
		TokenInformation:       emailAccount.TokenInformation,
		CreatedAt:              emailAccount.GetCreatedDate(),
	}

	return base, nil
}
