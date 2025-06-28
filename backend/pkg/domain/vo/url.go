package vo

import (
	"net/url"
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type URL struct {
	value string
}

func NewURL(rawURL string) (URL, error) {
	if rawURL == "" {
		return URL{}, nil
	}

	// Corrigido: tratar erro corretamente
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return URL{}, msgerror.AnErrInvalidURL
	}

	// Validação adicional para scheme e host
	if !isValidScheme(parsed.Scheme) || parsed.Host == "" {
		return URL{}, msgerror.AnErrInvalidURL
	}

	return URL{value: strings.TrimSuffix(rawURL, "/")}, nil
}

func isValidScheme(scheme string) bool {
	return scheme == "http" || scheme == "https"
}

func (u URL) String() string {
	return u.value
}

func (u URL) Equal(other URL) bool {
	return u.value == other.value
}

func (u URL) IsEmpty() bool {
	return u.value == ""
}
