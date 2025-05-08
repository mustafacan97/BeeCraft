package vo

import "platform/pkg/domain"

type OAuth2Credentials struct {
	domain.BaseValueObject
	clientID     string
	tenantID     string
	clientSecret string
}

func NewOAuth2Credentials(clientID, tenantID, clientSecret string) *OAuth2Credentials {
	return &OAuth2Credentials{
		clientID:     clientID,
		tenantID:     tenantID,
		clientSecret: clientSecret,
	}
}

func (e *OAuth2Credentials) GetAtomicValues() []interface{} {
	return []any{e.clientID, e.tenantID, e.clientSecret}
}

func (e *OAuth2Credentials) Credentials() (clientID, tenantID, clientSecret string) {
	return e.clientID, e.tenantID, e.clientSecret
}
