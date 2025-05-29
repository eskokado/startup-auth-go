package vo

import (
	"fmt"
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type Name string

// NewName cria um nome válido com limites opcionais (default: min=3, max=50)
func NewName(value string, min, max int) (Name, error) {
	if value == "" {
		return "", msgerror.AnErrEmptyName
	}

	trimmed := strings.TrimSpace(value)
	if min <= 0 {
		min = 3
	}
	if max <= 0 {
		max = 255
	}

	switch {
	case len(trimmed) < min:
		return "", fmt.Errorf("%w: mínimo %d caracteres", msgerror.AnErrNameTooShort, min)
	case len(trimmed) > max:
		return "", fmt.Errorf("%w: máximo %d caracteres", msgerror.AnErrNameTooLong, max)
	}

	return Name(trimmed), nil
}

func (n Name) String() string {
	return string(n)
}

func (n Name) Equal(other Name) bool {
	return n == other
}

func (n Name) IsEmpty() bool {
	return n == ""
}
