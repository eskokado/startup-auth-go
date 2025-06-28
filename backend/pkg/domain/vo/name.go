package vo

import (
	"fmt"
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

// Agora Name é uma struct com campo privado
type Name struct {
	value string
}

// NewName cria um nome válido com limites opcionais (default: min=3, max=255)
func NewName(value string, min, max int) (Name, error) {
	if value == "" {
		return Name{}, msgerror.AnErrEmptyName
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
		return Name{}, fmt.Errorf("%w: mínimo %d caracteres", msgerror.AnErrNameTooShort, min)
	case len(trimmed) > max:
		return Name{}, fmt.Errorf("%w: máximo %d caracteres", msgerror.AnErrNameTooLong, max)
	}

	return Name{value: trimmed}, nil
}

func (n Name) String() string {
	return n.value
}

func (n Name) Equal(other Name) bool {
	return n.value == other.value
}

func (n Name) IsEmpty() bool {
	return n.value == ""
}
