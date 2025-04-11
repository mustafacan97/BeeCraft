package valueobject

import (
	"errors"
	"platform/pkg/domain"
	"regexp"
)

var (
	emailRegex      = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	ErrEmailInvalid = errors.New("invalid email format")
)

type Email struct {
	domain.BaseValueObject
	value string
}

func NewEmail(value string) (Email, error) {
	if !emailRegex.MatchString(value) {
		return Email{}, ErrEmailInvalid
	}
	return Email{value: value}, nil
}

func (e Email) GetAtomicValues() []interface{} {
	return []any{e.value}
}

func (e Email) GetValue() string {
	return e.value
}
