package email_sender

import (
	"errors"
	"fmt"
	"net/smtp"
)

// OAuth2 implements the smtp.Auth interface using OAuth2 access tokens.
//
// username is the email address of the user being authenticated.
// accessToken is the OAuth2 token used to authenticate the user.
type OAuth2 struct {
	username    string
	accessToken string
}

// NewOAuth2Auth returns a new smtp.Auth implementation using the XOAUTH2 mechanism.
//
// It accepts a username (email address) and an OAuth2 access token to authenticate with
// SMTP servers like Gmail or Microsoft 365.
func NewOAuth2Auth(username, accessToken string) smtp.Auth {
	return &OAuth2{username: username, accessToken: accessToken}
}

// Start begins the XOAUTH2 authentication process.
//
// It returns the authentication mechanism name ("XOAUTH2"), the initial client response
// containing the formatted authentication string, and an error if something goes wrong.
func (a *OAuth2) Start(server *smtp.ServerInfo) (string, []byte, error) {
	authStr := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", a.username, a.accessToken)
	return "XOAUTH2", []byte(authStr), nil
}

// Next handles any additional challenges from the server.
//
// Since XOAUTH2 is a single-step authentication process, receiving a challenge
// indicates an unexpected server response and results in an error.
func (a *OAuth2) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}
