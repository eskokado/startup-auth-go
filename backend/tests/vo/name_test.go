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
		{"Limite Mínimo Padrão", "Ab", 0, 0, msgerror.AnErrNameTooShort},
		{"Limite Máximo Padrão", strings.Repeat("A", 256), 0, 0, msgerror.AnErrNameTooLong},
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
				return
			}

			// Verifica se espaços são removidos
			expectedTrimmed := strings.TrimSpace(tt.input)
			if got.String() != expectedTrimmed {
				t.Errorf("NewName() = %v, espera-se %v", got.String(), expectedTrimmed)
			}
		})
	}
}

func TestNameMethods(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		n, _ := vo.NewName("Alice", 3, 50)
		if n.String() != "Alice" {
			t.Errorf("String() = %v, want %v", n.String(), "Alice")
		}
	})

	t.Run("Equal()", func(t *testing.T) {
		n1, _ := vo.NewName("Alice", 3, 50)
		n2, _ := vo.NewName("Alice", 3, 50)
		n3, _ := vo.NewName("Bob", 3, 50)

		if !n1.Equal(n2) {
			t.Error("Equal() deve retornar true para nomes iguais")
		}
		if n1.Equal(n3) {
			t.Error("Equal() deve retornar false para nomes diferentes")
		}
	})

	t.Run("IsEmpty()", func(t *testing.T) {
		empty, _ := vo.NewName("", 0, 50)
		nonEmpty, _ := vo.NewName("Alice", 3, 50)

		if !empty.IsEmpty() {
			t.Error("IsEmpty() deve retornar true para nome vazio")
		}
		if nonEmpty.IsEmpty() {
			t.Error("IsEmpty() deve retornar false para nome não-vazio")
		}
	})
}
