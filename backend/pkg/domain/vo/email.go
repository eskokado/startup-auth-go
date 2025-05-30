package vo

import (
	"net/mail"
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type Email string

func NewEmail(value string) (Email, error) {
	if value == "" {
		return "", msgerror.AnErrEmptyEmail
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return "", msgerror.AnErrInvalidEmail
	}

	return Email(strings.ToLower(value)), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Equal(other Email) bool {
	return e == other
}

func (e Email) IsEmpty() bool {
	return e == ""
}
