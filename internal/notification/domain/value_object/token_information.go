package valueobject

import (
	"platform/pkg/domain"
	"time"
)

type TokenInformation struct {
	domain.BaseValueObject
	accessToken  string
	refreshToken string
	expireAt     time.Time
}

func NewTokenInformation(accessToken, refreshToken string, expireAt time.Time) *TokenInformation {
	return &TokenInformation{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expireAt:     expireAt,
	}
}

func (e *TokenInformation) GetTokenInformation() (string, string, time.Time) {
	return e.accessToken, e.refreshToken, e.expireAt
}

func (e *TokenInformation) SetTokenInformation(accessToken, refreshToken string, expireAt time.Time) {
	e.accessToken = accessToken
	e.refreshToken = refreshToken
	e.expireAt = expireAt
}

func (e *TokenInformation) GetAtomicValues() []interface{} {
	return []any{e.accessToken, e.refreshToken, e.expireAt}
}
