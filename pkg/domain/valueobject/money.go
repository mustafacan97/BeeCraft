package valueobject

import (
	"errors"
	"fmt"
	"platform/pkg/domain"
	"strings"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
)

type Money struct {
	domain.BaseValueObject
	Amount   float64
	Currency string
}

var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"TRY": true,
	"GBP": true,
}

func NewMoney(amount float64, currency string) (Money, error) {
	currency = strings.ToUpper(currency)
	if !validCurrencies[currency] {
		return Money{}, ErrInvalidCurrency
	}
	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m Money) GetAtomicValues() []interface{} {
	return []any{m.Amount, m.Currency}
}

func (m Money) String() string {
	return fmt.Sprintf("%f %s", m.Amount, m.Currency)
}
