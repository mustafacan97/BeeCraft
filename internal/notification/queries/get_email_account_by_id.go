package queries

import (
	"context"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	"platform/pkg/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

type GetEmailAccountByIDQuery struct {
	ID uuid.UUID
}

type GetEmailAccountByIDQueryResponse struct {
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
		EnableSSL:              emailAccount.IsSslEnabled(),
		TypeId:                 emailAccount.GetSmtpType(),
		TraditionalCredentials: emailAccount.TraditionalCredentials,
		OAuth2Credentials:      emailAccount.OAuth2Credentials,
		TokenInformation:       emailAccount.TokenInformation,
		CreatedAt:              emailAccount.GetCreatedDate(),
	}

	return base, nil
}
