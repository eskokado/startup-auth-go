package vo_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"golang.org/x/crypto/bcrypt"
)

func TestNewPasswordHash(t *testing.T) {
	// Gerar um hash bcrypt válido para testes
	validHash, _ := bcrypt.GenerateFromPassword([]byte("ValidPass123!"), bcrypt.DefaultCost)
	validHashStr := string(validHash)

	// Hash inválido com exatamente 60 caracteres
	invalidHash := "$2a$10$INVALIDHASHINVALIDHASHINVALIDHASHINVALID"

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"Senha válida", "SecurePass123!", nil},
		{"Senha curta", "short", msgerror.AnErrPasswordInvalid},
		{"Senha vazia", "", msgerror.AnErrPasswordInvalid},
		{"Caracteres insuficientes", "abc123", msgerror.AnErrPasswordInvalid},
		{"Hash bcrypt válido", validHashStr, nil},
		{"Hash bcrypt inválido", invalidHash, nil},
		{"Senha extremamente longa", strings.Repeat("a", 100), errors.New("bcrypt: password length exceeds 72 bytes")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := vo.NewPasswordHash(tt.password)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("NewPasswordHash() retornou erro nil, esperado %v", tt.wantErr)
				}

				if tt.wantErr.Error() != err.Error() {
					t.Errorf("NewPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewPasswordHash() unexpected error = %v", err)
			}
		})
	}
}

func TestPasswordHash_Verify(t *testing.T) {
	hash, err := vo.NewPasswordHash("ValidPass123!")
	if err != nil {
		t.Fatalf("Falha ao criar hash para teste: %v", err)
	}

	directHash, _ := vo.NewPasswordHash("DirectHash123!")

	tests := []struct {
		name     string
		hash     vo.PasswordHash
		password string
		expected bool
	}{
		{"Senha correta", hash, "ValidPass123!", true},
		{"Senha incorreta", hash, "WrongPass456@", false},
		{"Senha vazia", hash, "", false},
		{"Hash direto correto", directHash, "DirectHash123!", true},
		{"Hash direto incorreto", directHash, "WrongPass", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hash.Verify(tt.password)
			if result != tt.expected {
				t.Errorf("Verify() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPasswordHash_Methods(t *testing.T) {
	// Usar senha válida para os testes
	validPassword := "ValidPassword123!"

	t.Run("String()", func(t *testing.T) {
		ph, err := vo.NewPasswordHash(validPassword)
		if err != nil {
			t.Fatalf("Falha ao criar hash: %v", err)
		}
		if ph.String() == "" {
			t.Error("String() não deve retornar vazio para hash válido")
		}
	})

	t.Run("IsEmpty()", func(t *testing.T) {
		empty := vo.PasswordHash{}
		nonEmpty, err := vo.NewPasswordHash(validPassword)
		if err != nil {
			t.Fatalf("Falha ao criar hash: %v", err)
		}

		if !empty.IsEmpty() {
			t.Error("IsEmpty() deve retornar true para hash vazio")
		}
		if nonEmpty.IsEmpty() {
			t.Error("IsEmpty() deve retornar false para hash não-vazio")
		}
	})
}

func TestIsBcryptHash(t *testing.T) {
	// Gerar um hash bcrypt válido
	validHash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	validHashStr := string(validHash)

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Hash válido", validHashStr, true},
		{"Hash inválido curto", "short", false},
		{"Hash inválido formato", "$invalid$hash", false},
		{"String vazia", "", false},
		{"Outro tipo de hash", "sha256$abcdef", false},
		{"Comprimento 60 mas formato inválido", strings.Repeat("a", 60), false},
		{"Partes insuficientes", "part1$part2", false},
		{"Parte 0 não vazia", "prefix$2a$10$hash", false},
		{"Parte 1 não começa com 2", "part0$1a$10$hash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vo.IsBcryptHash(tt.input)
			if result != tt.want {
				t.Errorf("isBcryptHash() = %v, want %v", result, tt.want)
			}
		})
	}
}
