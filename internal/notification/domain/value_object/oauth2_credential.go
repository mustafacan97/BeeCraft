package valueobject

import (
	"errors"
	"platform/pkg/domain"
)

type OAuth2Credential struct {
	domain.BaseValueObject
	clientID     string
	tenantID     string
	clientSecret string
}

func NewOAuth2Credentials(clientID, tenantID, clientSecret string) (*OAuth2Credential, error) {
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("clientID and clientSecret required for OAuth2 credentials")
	}
	return &OAuth2Credential{
		clientID:     clientID,
		tenantID:     tenantID,
		clientSecret: clientSecret,
	}, nil
}

func (e *OAuth2Credential) GetCredentials() (clientID, tenantID, clientSecret string) {
	return e.clientID, e.tenantID, e.clientSecret
}

func (e *OAuth2Credential) SetCredentials(clientID, tenantID, clientSecret string) {
	e.clientID = clientID
	e.tenantID = tenantID
	e.clientSecret = clientSecret
}

func (e *OAuth2Credential) GetAtomicValues() []interface{} {
	return []any{e.clientID, e.tenantID, e.clientSecret}
}
