package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"platform/internal/notification/commands"
	"platform/internal/notification/domain"
	"platform/internal/notification/queries"
	"platform/internal/shared"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2CallbackRequest struct {
	Code         string `json:"code"`
	State        string `json:"state"`
	SessionState string `json:"session_state"`
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
}

type OAuth2CallbackResponse struct{}

type OAuth2CallbackHandler struct{}

func (h *OAuth2CallbackHandler) Handle(ctx context.Context, req *OAuth2CallbackRequest) (*baseHandler.Response[shared.HALResource], error) {
	if req.State == "" || req.Code == "" {
		return baseHandler.FailedResponse[shared.HALResource](errors.New("callback parameters are not valid")), nil
	}

	// STEP-1: Get encoded email account identifier from request
	decodedState, err := base64.StdEncoding.DecodeString(req.State)
	if err != nil {
		log.Printf("Failed to decode state: %v", err)
		return nil, fmt.Errorf("invalid state parameter")
	}

	// STEP-2: Convert []byte to string
	emailAccountID, err := uuid.FromBytes(decodedState)
	if err != nil {
		log.Printf("Invalid UUID format: %v", err)
		return nil, fmt.Errorf("invalid UUID in state parameter")
	}

	// STEP-3: Get email account from repository
	query := queries.GetEmailAccountByIDQuery{ID: emailAccountID}
	resp, err := mediator.Send[*queries.GetEmailAccountByIDQuery, *queries.GetEmailAccountByIDQueryResponse](ctx, &query)
	if err != nil {
		// TODO: log error
		return baseHandler.FailedResponse[shared.HALResource](errors.New("email account not found")), nil
	}
	if resp == nil {
		return baseHandler.FailedResponse[shared.HALResource](errors.New("email account not found")), nil
	}

	clientID, tenantID, clientSecret := resp.OAuth2Credentials.GetCredentials()
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/oauth2-callback",
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

	token, err := oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		return baseHandler.FailedResponse[shared.HALResource](errors.New("token is not valid")), nil
	}

	// STEP-4
	traditionalCredentials := resp.TraditionalCredentials
	oauth2Credentials := resp.OAuth2Credentials
	command := commands.UpdateEmailAccountCommand{
		ID:           emailAccountID,
		Email:        resp.Email.GetValue(),
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
		username, password := traditionalCredentials.GetCredentials()
		command.Username = username
		command.Password = password
	} else if oauth2Credentials != nil {
		clientID, tenantID, clientSecret := oauth2Credentials.GetCredentials()
		command.ClientID = clientID
		command.TenantID = tenantID
		command.ClientSecret = clientSecret
	}
	_, err = mediator.Send[*commands.UpdateEmailAccountCommand, *commands.UpdateEmailAccountCommandResponse](ctx, &command)
	if err != nil {
		return baseHandler.FailedResponse[shared.HALResource](err), nil
	}

	return baseHandler.SuccessResponse[shared.HALResource](nil), nil
}
