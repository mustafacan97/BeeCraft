package vo

import "platform/pkg/domain"

type TraditionalCredentials struct {
	domain.BaseValueObject
	username string
	password string
}

func NewTraditionalCredentials(username, password string) *TraditionalCredentials {
	return &TraditionalCredentials{
		username: username,
		password: password,
	}
}

func (e *TraditionalCredentials) GetAtomicValues() []interface{} {
	return []any{e.username, e.password}
}

func (e *TraditionalCredentials) Credentials() (string, string) {
	return e.username, e.password
}
