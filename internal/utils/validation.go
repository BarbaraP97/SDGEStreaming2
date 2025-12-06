// internal/utils/validation.go
package utils

import (
	"errors"
	"regexp"
	"strings"
)

// VALIDACIONES BÁSICAS (TUS FUNCIONES ORIGINALES)

// IsValidEmail checks if an email has a basic valid format.
func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}

// IsValidPassword checks if a password meets minimum requirements.
func IsValidPassword(password string) bool {
	return len(password) >= 6
}

// IsValidName checks if a name is valid (no numbers or special chars).
func IsValidName(name string) bool {
	if len(name) < 2 {
		return false
	}
	for _, r := range name {
		if (r >= '0' && r <= '9') || r == '@' || r == '#' || r == '$' {
			return false
		}
	}
	return true
}

// VALIDACIONES AVANZADAS (MEJORA PROFESIONAL, COMPLEMENTA LO EXISTENTE)

// Regex para validación más estricta de email (opcional, no rompe tu lógica actual)
var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

// ValidateEmail returns an error if email is invalid.
func ValidateEmail(email string) error {
	if IsEmpty(email) {
		return errors.New("el correo electrónico es obligatorio")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("el correo electrónico no es válido")
	}
	return nil
}

// ValidatePassword returns an error if password is invalid.
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("la contraseña debe tener al menos 6 caracteres")
	}
	return nil
}
