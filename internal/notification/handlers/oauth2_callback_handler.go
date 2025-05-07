package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"platform/internal/notification/domain"
	"platform/internal/notification/mediatr/commands"
	event_notification "platform/internal/notification/mediatr/notifications"
	"platform/internal/notification/mediatr/queries"
	"platform/internal/shared"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2CallbackRequest struct {
	Code             string `reqHeader:"-" params:"-" query:"code" json:"-" validate:"required"`
	Error            string `reqHeader:"-" params:"-" query:"error" json:"-"`
	ErrorDescription string `reqHeader:"-" params:"-" query:"error_description" json:"-"`
	ErrorUri         string `reqHeader:"-" params:"-" query:"error_uri" json:"-"`
	State            string `reqHeader:"-" params:"-" query:"state" json:"-" validate:"required"`
}

type OAuth2CallbackResponse struct {
}

type OAuth2CallbackHandler struct{}

func (h *OAuth2CallbackHandler) Handle(ctx context.Context, req *OAuth2CallbackRequest) (*baseHandler.Response[OAuth2CallbackResponse], error) {
	if req.Error != "" {
		return nil, fmt.Errorf("an error occurred on oauth2 callback with description: %s", req.ErrorDescription)
	}

	// STEP-1: Get encoded email account identifier from request
	decodedState, err := base64.StdEncoding.DecodeString(req.State)
	if err != nil {
		log.Printf("Failed to decode state: %v", err)
		return nil, fmt.Errorf("invalid state parameter")
	}

	// STEP-2: Unmarshal into structured state object
	var stateObj struct {
		ProjectID string `json:"project_id"`
		Email     string `json:"email"`
	}
	if err := json.Unmarshal(decodedState, &stateObj); err != nil {
		log.Printf("Failed to unmarshal state: %v", err)
		return nil, fmt.Errorf("invalid state content")
	}

	// STEP-3: Check if variables are valid which are coming from state
	email := stateObj.Email
	projectID, err := uuid.Parse(stateObj.ProjectID)
	if err != nil {
		return nil, shared.ErrInvalidContext
	}

	// STEP-4: Save project identifier to users' context
	ctx = context.WithValue(ctx, shared.ProjectIDContextKey, projectID)

	// STEP-5: Get email account from repository
	query := queries.GetEmailAccountByEmailQuery{Email: email}
	resp, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return baseHandler.NotFoundResponse[OAuth2CallbackResponse](), nil
	}

	// STEP-6: Get email account credentials
	clientID, tenantID, clientSecret := resp.OAuth2Credentials.Credentials()
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/v1/notification/email-accounts/oauth2-callback",
	}

	if resp.TypeId == domain.GmailOAuth2 {
		oauth2Config.Scopes = []string{"https://mail.google.com/"}
		oauth2Config.Endpoint = google.Endpoint
	} else if resp.TypeId == domain.MicrosoftOAuth2 {
		oauth2Config.Scopes = []string{"https://outlook.office365.com/SMTP.Send", "offline_access"}
		oauth2Config.Endpoint = oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		}
	}

	// STEP-7: Create a token
	token, err := oauth2Config.Exchange(ctx, req.Code)
	if err != nil {
		zap.L().Error("Failed to exchange code for token", zap.Error(err))
		return nil, err
	}

	// STEP-8: Create update email account command
	traditionalCredentials := resp.TraditionalCredentials
	oauth2Credentials := resp.OAuth2Credentials
	command := commands.UpdateEmailAccountCommand{
		Email:        resp.Email.Value(),
		DisplayName:  resp.DisplayName,
		Host:         resp.Host,
		Port:         resp.Port,
		EnableSSL:    resp.EnableSSL,
		TypeID:       resp.TypeId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpireAt:     token.Expiry,
	}
	if traditionalCredentials != nil {
		username, password := traditionalCredentials.Credentials()
		command.Username = username
		command.Password = password
	} else if oauth2Credentials != nil {
		clientID, tenantID, clientSecret := oauth2Credentials.Credentials()
		command.ClientID = clientID
		command.TenantID = tenantID
		command.ClientSecret = clientSecret
	}
	_, err = mediator.Send[*commands.UpdateEmailAccountCommand, *commands.UpdateEmailAccountCommandResponse](ctx, &command)
	if err != nil {
		return nil, err
	}

	// STEP-9: Publish email account update notification
	notification := event_notification.NewEmailAccountUpdatedEvent(projectID, email)
	mediator.Publish(ctx, &notification)

	// STEP-10: Return hateoas links to user
	respData := OAuth2CallbackResponse{}
	response := baseHandler.SuccessResponse(&respData)
	response.Links = hateoasLinksForOAuth2Callback(email)
	return response, nil
}

func hateoasLinksForOAuth2Callback(email string) shared.HALLinks {
	return shared.HALLinks{
		"delete": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "DELETE",
			Title:  "Delete this email account",
		},
		"list": {
			Href:   "/v1/notification/email-accounts?p=1&ps=10",
			Method: "GET",
			Title:  "List all emails on the first page",
		},
		"self": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "GET",
			Title:  "View this email account",
		},
		"update": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "PUT",
			Title:  "Update this email account",
		},
	}
}
