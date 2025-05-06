package vo

import (
	"errors"
	"platform/pkg/domain"
	"regexp"
)

var (
	phoneRegex      = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	ErrPhoneInvalid = errors.New("invalid phone number format")
)

type PhoneNumber struct {
	domain.BaseValueObject
	value string
}

func NewPhoneNumber(value string) (PhoneNumber, error) {
	if !phoneRegex.MatchString(value) {
		return PhoneNumber{}, ErrPhoneInvalid
	}
	return PhoneNumber{value: value}, nil
}

func (p PhoneNumber) GetAtomicValues() []any {
	return []any{p.value}
}
