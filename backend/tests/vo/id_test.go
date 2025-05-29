package vo

import (
	"errors"
	"testing"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/google/uuid"
)

func TestNewID(t *testing.T) {
	id := vo.NewID()

	if id.String() == "" {
		t.Errorf("NewID() retornou um ID vazio")
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"Válido", uuid.New().String(), nil},
		{"Vazio", "", msgerror.AnErrEmptyID},
		{"Inválido Formato", "invalid-uuid", msgerror.AnErrInvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := vo.ParseID(tt.input)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("ParseID() erro = %v, esperado %v", err, tt.wantErr)
				}
				return
			}

			if id.String() != tt.input {
				t.Errorf("ParseID() = %v, esperado %v", id.String(), tt.input)
			}
		})
	}
}

func TestIDEqual(t *testing.T) {
	id1 := vo.NewID()
	id2 := vo.NewID()

	if id1.Equal(id2) {
		t.Errorf("IDEqual() falhou: IDs diferentes não devem ser iguais")
	}

	if !id1.Equal(id1) {
		t.Errorf("IDEqual() falhou: O mesmo ID deve ser igual a si mesmo")
	}
}
