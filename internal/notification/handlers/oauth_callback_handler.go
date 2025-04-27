package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"platform/internal/notification/domain"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/internal/notification/repositories"
	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthCallbackRequest struct {
	Code         string `json:"code"`
	State        string `json:"state"`
	SessionState string `json:"session_state"`
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
}

type OAuthCallbackResponse struct {
}

type OAuthCallbackHandler struct {
	emailAccountRepository repositories.EmailAccountRepository
}

func NewOAuthCallbackHandler(emailAccountRepository *repositories.EmailAccountRepository) *OAuthCallbackHandler {
	return &OAuthCallbackHandler{
		emailAccountRepository: *emailAccountRepository,
	}
}

func (h *OAuthCallbackHandler) Handle(ctx context.Context, req *OAuthCallbackRequest) (*baseHandler.Response[string], error) {
	if req.State == "" {
		return baseHandler.FailedResponse[string](errors.New("Email account not found.")), nil
	}

	if req.Code == "" {
		return baseHandler.FailedResponse[string](errors.New("Authorization code not found.")), nil
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
	emailAccount, err := h.emailAccountRepository.GetByID(ctx, emailAccountID)
	if err != nil {
		return baseHandler.FailedResponse[string](errors.New("Email account not found.")), nil
	}

	values := emailAccount.OAuth2Credentials.GetAtomicValues()
	clientID := values[0].(string)
	tenantID := values[1].(string)
	clientSecret := values[2].(string)
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("%s://%s/oauth-callback", "http", req.Host),
	}

	typeID := emailAccount.GetSmtpType()
	if typeID == domain.GmailOAuth2 {
		oauth2Config.Scopes = []string{"https://mail.google.com/"}
		oauth2Config.Endpoint = google.Endpoint
	} else if typeID == domain.MicrosoftOAuth2 {
		oauth2Config.Scopes = []string{"https://outlook.office365.com/SMTP.Send", "offline_access"}
		oauth2Config.Endpoint = oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		}
	}

	token, err := oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		return baseHandler.FailedResponse[string](errors.New("Token is not valid.")), nil
	}

	// STEP-4
	emailAccount.TokenInformation = internalValueObject.NewTokenInformation(token.AccessToken, token.RefreshToken, token.Expiry)
	err = h.emailAccountRepository.Update(context.Background(), emailAccount)
	if err != nil {
		return baseHandler.FailedResponse[string](errors.New("An error occurred when saving access token.")), nil
	}

	return baseHandler.SuccessResponse[string](nil), nil
}
