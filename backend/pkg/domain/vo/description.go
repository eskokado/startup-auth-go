package vo

import (
	"fmt"
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type Description struct {
	value string
}

// NewDescription cria uma descrição válida com limites opcionais (default: min=3, max=50)
func NewDescription(value string, min, max int) (Description, error) {
	if value == "" {
		return Description{}, msgerror.AnErrEmptyDescription
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
		return Description{}, fmt.Errorf("%w: mínimo %d caracteres", msgerror.AnErrTooShort, min)
	case len(trimmed) > max:
		return Description{}, fmt.Errorf("%w: máximo %d caracteres", msgerror.AnErrTooLong, max)
	}

	return Description{value: trimmed}, nil
}

func (d Description) String() string {
	return d.value
}

func (d Description) Equal(other Description) bool {
	return d.value == other.value
}

func (d Description) IsEmpty() bool {
	return d.value == ""
}
