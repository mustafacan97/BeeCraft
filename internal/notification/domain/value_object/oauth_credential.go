package valueobject

import "platform/pkg/domain"

type OAuthCredential struct {
	domain.BaseValueObject
	clientID     string
	clientSecret string
	tenantID     string
}

func NewOAuthCredential(clientID, clientSecret, tenantID string) *OAuthCredential {
	return &OAuthCredential{
		clientID:     clientID,
		clientSecret: clientSecret,
		tenantID:     tenantID,
	}
}

func (e *OAuthCredential) GetCredentials() (string, string, string) {
	return e.clientID, e.clientSecret, e.tenantID
}

func (e *OAuthCredential) GetAtomicValues() []interface{} {
	return []any{e.clientID, e.clientSecret, e.tenantID}
}
