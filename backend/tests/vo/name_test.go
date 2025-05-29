package vo

import (
	"errors"
	"strings"
	"testing"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

func TestNewName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		min     int
		max     int
		wantErr error
	}{
		{"Válido", "John Doe", 3, 50, nil},
		{"Com Espaços", "  John  ", 3, 50, nil},
		{"Vazio", "", 3, 50, msgerror.AnErrEmptyName},
		{"Muito Curto", "A", 3, 50, msgerror.AnErrNameTooShort},
		{"Muito Longo", strings.Repeat("A", 51), 3, 50, msgerror.AnErrNameTooLong},
		{"Caracteres Especiais", "João Çãô", 3, 50, nil},

		// Novos casos para limites padrão
		{"Limite Mínimo Padrão", "Ab", 0, 0, msgerror.AnErrNameTooShort},                    // min padrão=3
		{"Limite Máximo Padrão", strings.Repeat("A", 256), 0, 0, msgerror.AnErrNameTooLong}, // max padrão=255
		{"Limites Personalizados", "AB", 2, 100, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewName(tt.input, tt.min, tt.max)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewName() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewName() unexpected error = %v", err)
			}

			// Verifica se espaços são removidos
			if strings.TrimSpace(tt.input) != string(got) {
				t.Errorf("NewName() = %v, espera-se espaços removidos", got)
			}
		})
	}
}

func TestNameMethods(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		n := vo.Name("Alice")
		if n.String() != "Alice" {
			t.Errorf("String() = %v, want %v", n.String(), "Alice")
		}
	})

	t.Run("Equal()", func(t *testing.T) {
		n1 := vo.Name("Alice")
		n2 := vo.Name("Alice")
		n3 := vo.Name("Bob")

		if !n1.Equal(n2) {
			t.Error("Equal() deve retornar true para nomes iguais")
		}
		if n1.Equal(n3) {
			t.Error("Equal() deve retornar false para nomes diferentes")
		}
	})

	t.Run("IsEmpty()", func(t *testing.T) {
		empty := vo.Name("")
		nonEmpty := vo.Name("Alice")

		if !empty.IsEmpty() {
			t.Error("IsEmpty() deve retornar true para nome vazio")
		}
		if nonEmpty.IsEmpty() {
			t.Error("IsEmpty() deve retornar false para nome não-vazio")
		}
	})
}
