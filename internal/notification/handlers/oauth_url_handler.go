package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"platform/internal/notification/domain"
	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthUrlRequest struct {
	EmailAccountID uuid.UUID
	ClientID       string
	ClientSecret   string
	TenantID       string
	RedirectURL    string
	TypeID         int
}

type OAuthUrlResponse struct {
}

type OAuthUrlHandler struct {
}

func NewOAuthUrlHandler() *OAuthUrlHandler {
	return &OAuthUrlHandler{}
}

func (h *OAuthUrlHandler) Handle(ctx context.Context, req *OAuthUrlRequest) (*baseHandler.Response[string], error) {
	oauth2Config := &oauth2.Config{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURL:  req.RedirectURL,
	}

	if req.TypeID == domain.GmailOAuth2 {
		oauth2Config.Scopes = []string{"https://mail.google.com/"}
		oauth2Config.Endpoint = google.Endpoint
	} else if req.TypeID == domain.MicrosoftOAuth2 {
		oauth2Config.Scopes = []string{"https://outlook.office365.com/SMTP.Send", "offline_access"}
		oauth2Config.Endpoint = oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", req.TenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", req.TenantID),
		}
	}

	encodedState := base64.StdEncoding.EncodeToString([]byte(req.EmailAccountID.String()))
	url := oauth2Config.AuthCodeURL(encodedState, oauth2.AccessTypeOffline)
	return baseHandler.SuccessResponse(&url), nil
}
