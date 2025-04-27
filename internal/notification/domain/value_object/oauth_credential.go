package valueobject

import "platform/pkg/domain"

type OAuth2Credential struct {
	domain.BaseValueObject
	clientID     string
	tenantID     string
	clientSecret string
}

func NewOAuth2Credentials(clientID, tenantID, clientSecret string) *OAuth2Credential {
	return &OAuth2Credential{
		clientID:     clientID,
		tenantID:     tenantID,
		clientSecret: clientSecret,
	}
}

func (e *OAuth2Credential) GetCredentials() (clientID, tenantID, clientSecret string) {
	return e.clientID, e.tenantID, e.clientSecret
}

func (e *OAuth2Credential) GetAtomicValues() []interface{} {
	return []any{e.clientID, e.tenantID, e.clientSecret}
}
