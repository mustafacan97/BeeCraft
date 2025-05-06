package validators

import (
	"net"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func HostnameValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Check if the value is a valid IP address
	if net.ParseIP(value) != nil {
		return true
	}

	// Check if the value is a valid hostname
	hostnameRegex := `^(?i)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z0-9]{2,})$`
	re := regexp.MustCompile(hostnameRegex)
	return re.MatchString(value)
}
