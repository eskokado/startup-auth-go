package vo

import (
	"testing"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

func TestNewURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"Válido HTTP", "http://example.com", nil},
		{"Válido HTTPS", "https://example.com/path?query=param", nil},
		{"Inválido Scheme", "ftp://example.com", msgerror.AnErrInvalidURL},
		{"Inválido Formato", "htp://invalid", msgerror.AnErrInvalidURL},
		{"Sem Host", "http://", msgerror.AnErrInvalidURL},
		{"Vazio", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := vo.NewURL(tt.input)

			if tt.wantErr != nil {
				// Comparação direta de erros em vez de errors.Is
				if err == nil || err != tt.wantErr {
					t.Errorf("NewURL() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewURL() unexpected error = %v", err)
			}
		})
	}
}

func TestURLMethods(t *testing.T) {
	t.Run("Equal()", func(t *testing.T) {
		u1, _ := vo.NewURL("https://example.com")
		u2, _ := vo.NewURL("https://example.com")
		u3, _ := vo.NewURL("https://different.com")

		if !u1.Equal(u2) {
			t.Error("Equal() deve retornar true para URLs iguais")
		}

		if u1.Equal(u3) {
			t.Error("Equal() deve retornar false para URLs diferentes")
		}
	})

	t.Run("String()", func(t *testing.T) {
		u, _ := vo.NewURL("https://example.com")
		if u.String() != "https://example.com" {
			t.Errorf("String() = %v, want %v", u.String(), "https://example.com")
		}
	})

	t.Run("IsEmpty()", func(t *testing.T) {
		empty, _ := vo.NewURL("")
		nonEmpty, _ := vo.NewURL("https://example.com")

		if !empty.IsEmpty() {
			t.Error("IsEmpty() deve retornar true para URL vazia")
		}
		if nonEmpty.IsEmpty() {
			t.Error("IsEmpty() deve retornar false para URL não-vazia")
		}
	})
}
