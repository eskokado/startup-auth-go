package vo

import (
	"net/mail"
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

// Email agora é uma struct com campo privado
type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, msgerror.AnErrEmptyEmail
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return Email{}, msgerror.AnErrInvalidEmail
	}

	return Email{value: strings.ToLower(value)}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equal(other Email) bool {
	return e.value == other.value
}

func (e Email) IsEmpty() bool {
	return e.value == ""
}
