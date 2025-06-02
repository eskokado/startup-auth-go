package vo

import (
	"errors"
	"testing"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

func TestNewDescription(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		min     int
		max     int
		wantErr error
	}{
		{"Válido", "Descrição válida", 3, 50, nil},
		{"Vazio", "", 3, 50, msgerror.AnErrEmptyDescription},
		{"Muito Curto", "ab", 3, 50, msgerror.AnErrTooShort},
		{"Muito Longo", "Lorem ipsum dolor sit amet consectetur adipiscing elit...", 3, 50, msgerror.AnErrTooLong},
		{"Limite Mínimo", "abc", 3, 50, nil},
		{"Limite Máximo", "Lorem ipsum dolor sit amet consectetur adipiscin", 3, 50, nil},
		{"Limite Mínimo Zero", "abc", 0, 50, nil},
		{"Limite Máximo Zero", "Lorem ipsum dolor sit amet consectetur adipiscin", 3, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := vo.NewDescription(tt.input, tt.min, tt.max)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDescriptionMethods(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		n, _ := vo.NewDescription("loren ipsum dolor sit amet", 3, 50)
		if n.String() != "loren ipsum dolor sit amet" {
			t.Errorf("String() = %v, want %v", n.String(), "loren ipsum dolor sit amet")
		}
	})

	t.Run("Equal()", func(t *testing.T) {
		d1, _ := vo.NewDescription("test", 3, 50)
		d2, _ := vo.NewDescription("test", 3, 50)
		d3, _ := vo.NewDescription("different", 3, 50)

		if !d1.Equal(d2) {
			t.Error("Equal() deve retornar true para descrições iguais")
		}

		if d1.Equal(d3) {
			t.Error("Equal() deve retornar false para descrições diferentes")
		}
	})

	t.Run("IsEmpty()", func(t *testing.T) {
		empty, _ := vo.NewDescription("", 0, 50)
		nonEmpty, _ := vo.NewDescription("test", 3, 50)

		if !empty.IsEmpty() {
			t.Error("IsEmpty() deve retornar true para descrição vazia")
		}

		if nonEmpty.IsEmpty() {
			t.Error("IsEmpty() deve retornar false para descrição não-vazia")
		}
	})
}
