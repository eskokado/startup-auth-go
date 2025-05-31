package vo

import (
	"strings"
	"testing"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"Válido", "user@example.com", nil},
		{"Válido Maiúsculas", "USER@EXAMPLE.COM", nil},
		{"Inválido Formato", "invalid.email", msgerror.AnErrInvalidEmail},
		{"Vazio", "", msgerror.AnErrEmptyEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vo.NewEmail(tt.input)

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewEmail() unexpected error = %v", err)
				return
			}

			expected := strings.ToLower(tt.input)
			if email.String() != expected {
				t.Errorf("NewEmail() = %v, want %v", email.String(), expected)
			}
		})
	}
}

func TestEmailMethods(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		n, _ := vo.NewEmail("alice@email.com")
		if n.String() != "alice@email.com" {
			t.Errorf("String() = %v, want %v", n.String(), "alice@email.com")
		}
	})

	t.Run("Equal()", func(t *testing.T) {
		n1, _ := vo.NewEmail("alice@email.com")
		n2, _ := vo.NewEmail("alice@email.com")
		n3, _ := vo.NewEmail("bob@email.com")

		if !n1.Equal(n2) {
			t.Error("Equal() deve retornar true para emails iguais")
		}
		if n1.Equal(n3) {
			t.Error("Equal() deve retornar false para emails diferentes")
		}
	})

	t.Run("IsEmpty()", func(t *testing.T) {
		empty, _ := vo.NewEmail("")
		nonEmpty, _ := vo.NewEmail("alice@email.com")

		if !empty.IsEmpty() {
			t.Error("IsEmpty() deve retornar true para email vazio")
		}
		if nonEmpty.IsEmpty() {
			t.Error("IsEmpty() deve retornar false para email não-vazio")
		}
	})
}
