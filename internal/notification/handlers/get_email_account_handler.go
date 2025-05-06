package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"platform/internal/notification/domain"
	"platform/internal/notification/mediatr/queries"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	OAuth2Url    string `json:"oauth2_url"`
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

	if resp.OAuth2Credentials != nil {
		response.OAuth2Url = getOAuth2Url(req.ID, response.ClientID, response.TenantID, response.ClientSecret)
	}

	return baseHandler.SuccessResponse(&response), nil
}

func getOAuth2Url(emailAccountID uuid.UUID, clientID, tenantID, clientSecret string) string {
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/oauth2-callback",
	}

	if tenantID != "" {
		oauth2Config.Scopes = []string{"https://outlook.office365.com/SMTP.Send", "offline_access"}
		oauth2Config.Endpoint = oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		}
	} else {
		oauth2Config.Scopes = []string{"https://mail.google.com/"}
		oauth2Config.Endpoint = google.Endpoint
	}

	encodedState := base64.StdEncoding.EncodeToString([]byte(emailAccountID.String()))
	return oauth2Config.AuthCodeURL(encodedState, oauth2.AccessTypeOffline)
}
