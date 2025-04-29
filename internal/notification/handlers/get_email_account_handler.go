package handlers

import (
	"context"
	"platform/internal/notification/domain"
	"platform/internal/notification/queries"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"

	"github.com/google/uuid"
)

type GetEmailAccountRequest struct {
	ID uuid.UUID `params:"id" json:"-" validate:"required,uuid4"`
}

type GetEmailAccountResponse struct {
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	EnableSSL    bool   `json:"enable_ssl"`
	TypeId       int    `json:"type_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	TenantID     string `json:"tenant_id"`
	ClientSecret string `json:"client_secret"`
}

type GetEmailAccountHandler struct{}

func (h *GetEmailAccountHandler) Handle(ctx context.Context, req *GetEmailAccountRequest) (*baseHandler.Response[GetEmailAccountResponse], error) {
	query := queries.GetEmailAccountByIDQuery{ID: req.ID}
	resp, err := mediator.Send[*queries.GetEmailAccountByIDQuery, *queries.GetEmailAccountByIDQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return baseHandler.NotFoundResponse[GetEmailAccountResponse](), nil
	}

	response := GetEmailAccountResponse{
		Email:       resp.Email.GetValue(),
		DisplayName: resp.DisplayName,
		Host:        resp.Host,
		Port:        resp.Port,
		EnableSSL:   resp.EnableSSL,
		TypeId:      resp.TypeId,
	}

	switch resp.TypeId {
	case domain.Login:
		username, password := resp.TraditionalCredentials.GetCredentials()
		response.Username = username
		response.Password = password
	case domain.GmailOAuth2, domain.MicrosoftOAuth2:
		clientID, tenantID, clientSecret := resp.OAuth2Credentials.GetCredentials()
		response.ClientID = clientID
		response.TenantID = tenantID
		response.ClientSecret = clientSecret
	}

	return baseHandler.SuccessResponse(&response), nil
}
