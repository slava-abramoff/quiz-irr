package validator

import (
	"net/mail"
	"strings"
)

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if addr.Address != email {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) < 2 || !strings.Contains(parts[1], ".") {
		return false
	}

	return true
}
