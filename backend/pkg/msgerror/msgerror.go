package msgerror

import (
	"errors"
	"fmt"
)

type ValidationErrors struct {
	FieldErrors map[string]string // Campo -> mensagem
}

var (
	AnErrInvalidCredentials = errors.New("invalid credentials")
	AnErrUserNotFound       = errors.New("user not found")
	AnErrUserExists         = errors.New("user already exists")
	AnErrWeakPassword       = errors.New("password does not meet security requirements")
	AnErrInvalidUser        = errors.New("invalid user")
	AnErrEmptyEmail         = errors.New("email cannot be empty")
	AnErrNotFound           = errors.New("not found")
	AnErrEmptyDescription   = errors.New("description cannot be empty")
	AnErrTooShort           = errors.New("description too short")
	AnErrTooLong            = errors.New("description too long")
	AnErrInvalidEmail       = errors.New("invalid email format")
	AnErrEmptyID            = errors.New("empty ID")
	AnErrInvalidID          = errors.New("invalid ID format")
	AnErrInvalidName        = errors.New("invalid name")
	AnErrEmptyName          = errors.New("name cannot be empty")
	AnErrNameDifferent      = errors.New("new name must be different")
	AnErrPasswordInvalid    = errors.New("password must be at least 8 characters")
	AnErrEmptyPassword      = errors.New("expected error for empty password")
	AnErrInvalidURL         = errors.New("invalid URL format")
	AnErrEmptyURL           = errors.New("URL cannot be empty")
	AnErrCreateUser         = errors.New("failed to create user ")
	AnErrNoSavedUser        = errors.New("user not saved correctly")
	AnErrNameTooShort       = errors.New("name too short")
	AnErrNameTooLong        = errors.New("name too long")
	AnErrInvalidPassword    = errors.New("invalid password")
	AnErrInvalidToken       = errors.New("invalid token")
	AnErrExpiredToken       = errors.New("expired token")
	AnErrSendMessageByEmail = errors.New("error send message by email")
)

func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func (v ValidationErrors) Error() string {
	return "validation failed with multiple errors"
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		FieldErrors: make(map[string]string),
	}
}

func (v *ValidationErrors) Add(field, message string) {
	v.FieldErrors[field] = message
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.FieldErrors) > 0
}
