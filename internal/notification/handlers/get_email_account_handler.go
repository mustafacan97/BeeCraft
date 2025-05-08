package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"platform/internal/notification/domain"
	"platform/internal/notification/mediatr/queries"
	"platform/internal/shared"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GetEmailAccountRequest struct {
	ProjectID uuid.UUID `reqHeader:"X-Project-ID" params:"-" json:"-" validate:"required,uuid"`
	Email     string    `reqHeader:"-" params:"email" json:"-" validate:"required,email"`
}

type GetEmailAccountResponse struct {
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	EnableSSL    bool   `json:"enable_ssl"`
	TypeId       int    `json:"type_id"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	OAuth2Url    string `json:"oauth2_url,omitempty"`
}

type GetEmailAccountHandler struct{}

func (h *GetEmailAccountHandler) Handle(ctx context.Context, req *GetEmailAccountRequest) (*baseHandler.Response[GetEmailAccountResponse], error) {
	// STEP-1: Get the email account
	query := queries.GetEmailAccountByEmailQuery{Email: req.Email}
	resp, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return baseHandler.NotFoundResponse[GetEmailAccountResponse](), nil
	}

	// STEP-2: Create response struct
	data := GetEmailAccountResponse{
		Email:       resp.Email.Value(),
		DisplayName: resp.DisplayName,
		Host:        resp.Host,
		Port:        resp.Port,
		EnableSSL:   resp.EnableSSL,
		TypeId:      resp.TypeId,
	}

	// STEP-3: Get the related credentials
	if resp.TypeId == domain.Login {
		username, password := resp.TraditionalCredentials.Credentials()
		data.Username = username
		data.Password = password
	} else {
		clientID, tenantID, clientSecret := resp.OAuth2Credentials.Credentials()
		data.ClientID = clientID
		data.TenantID = tenantID
		data.ClientSecret = clientSecret
	}

	// STEP-4: Get OAut2 URL
	if resp.OAuth2Credentials != nil {
		data.OAuth2Url = getOAuth2Url(req.ProjectID, req.Email, data.ClientID, data.TenantID, data.ClientSecret)
	}

	// STEP-4: Returns hateoas links to user
	response := baseHandler.SuccessResponse(&data)
	response.Links = hateoasLinksForGet(req.Email)
	return response, nil
}

func getOAuth2Url(projectID uuid.UUID, email, clientID, tenantID, clientSecret string) string {
	if clientID == "" || tenantID == "" || clientSecret == "" {
		return ""
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/v1/notification/email-accounts/oauth2-callback",
	}

	if tenantID != "" {
		oauth2Config.Scopes = []string{
			"https://outlook.office365.com/SMTP.Send",
			"offline_access",
		}
		oauth2Config.Endpoint = oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		}
	} else {
		oauth2Config.Scopes = []string{"https://mail.google.com/"}
		oauth2Config.Endpoint = google.Endpoint
	}

	stateObj := map[string]string{
		"project_id": projectID.String(),
		"email":      email,
	}
	stateBytes, _ := json.Marshal(stateObj)
	encodedState := base64.URLEncoding.EncodeToString(stateBytes)

	return oauth2Config.AuthCodeURL(encodedState, oauth2.AccessTypeOffline)
}

func hateoasLinksForGet(email string) shared.HALLinks {
	return shared.HALLinks{
		"delete": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "DELETE",
			Title:  "Delete this email account",
		},
		"list": {
			Href:   "/v1/notification/email-accounts?p=0&ps=1",
			Method: "GET",
			Title:  "List all emails on the first page",
		},
		"update": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "PUT",
			Title:  "Update this email account",
		},
	}
}
