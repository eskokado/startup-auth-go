package entity_test

import (
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		inputEmail  string
		inputPass   string
		inputImgURL string
		wantErr     error
	}{
		{
			name:        "Usuário válido",
			inputName:   "John Doe",
			inputEmail:  "john@example.com",
			inputPass:   "SecurePass123!",
			inputImgURL: "https://example.com/avatar.jpg",
			wantErr:     nil,
		},
		{
			name:        "Email inválido",
			inputName:   "John Doe",
			inputEmail:  "invalid-email",
			inputPass:   "SecurePass123!",
			inputImgURL: "https://example.com/avatar.jpg",
			wantErr:     msgerror.AnErrInvalidUser,
		},
		{
			name:        "Senha fraca",
			inputName:   "John Doe",
			inputEmail:  "john@example.com",
			inputPass:   "short",
			inputImgURL: "https://example.com/avatar.jpg",
			wantErr:     msgerror.AnErrInvalidUser,
		},
		{
			name:        "URL de imagem inválida",
			inputName:   "John Doe",
			inputEmail:  "john@example.com",
			inputPass:   "SecurePass123!",
			inputImgURL: "htp://invalid-url",
			wantErr:     msgerror.AnErrInvalidUser,
		},
		{
			name:        "Campos obrigatórios vazios",
			inputName:   "",
			inputEmail:  "",
			inputPass:   "",
			inputImgURL: "",
			wantErr:     msgerror.AnErrInvalidUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.CreateUser(
				tt.inputName,
				tt.inputEmail,
				tt.inputPass,
				tt.inputImgURL,
			)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateUser() unexpected error = %v", err)
			}
		})
	}
}

func TestUser_Equal(t *testing.T) {
	user1, _ := entity.CreateUser(
		"Alice",
		"alice@example.com",
		"StrongPass123!",
		"",
	)

	user2, _ := entity.CreateUser(
		"Bob",
		"bob@example.com",
		"AnotherPass123!",
		"",
	)

	tests := []struct {
		name     string
		userA    *entity.User
		userB    *entity.User
		expected bool
	}{
		{"Mesmo usuário", user1, user1, true},
		{"IDs diferentes", user1, user2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userA.Equal(tt.userB)
			if result != tt.expected {
				t.Errorf("Equal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test User",
		"test@example.com",
		"CorrectPass123!",
		"",
	)

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Senha correta", "CorrectPass123!", true},
		{"Senha incorreta", "WrongPass456@", false},
		{"Senha vazia", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.VerifyPassword(tt.password)
			if result != tt.expected {
				t.Errorf("VerifyPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewUser_Validations(t *testing.T) {
	validID := vo.NewID()
	validName, _ := vo.NewName("Test", 3, 50)
	validEmail, _ := vo.NewEmail("test@example.com")
	emptyEmail, _ := vo.NewEmail("")
	emptyHash, _ := vo.NewPasswordHash("")
	validHash, _ := vo.NewPasswordHash("ValidPass123!")

	tests := []struct {
		name         string
		email        vo.Email
		passwordHash vo.PasswordHash
		wantErr      error
	}{
		{"Email vazio", emptyEmail, validHash, msgerror.AnErrEmptyEmail},
		{"Hash vazio", validEmail, emptyHash, msgerror.AnErrWeakPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.NewUser(
				validID,
				validName,
				tt.email,
				tt.passwordHash,
				nil,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_WithName_Error(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)
	emptyName, _ := vo.NewName("", 0, 50)
	_, err := user.WithName(vo.Name(emptyName)) // Nome vazio
	if !errors.Is(err, msgerror.AnErrInvalidName) {
		t.Errorf("Esperado erro %v, obtido %v", msgerror.AnErrInvalidName, err)
	}
}

func TestUser_WithPasswordHash_Error(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)
	emptyHash, _ := vo.NewPasswordHash("")

	_, err := user.WithPasswordHash(emptyHash) // Hash vazio
	if !errors.Is(err, msgerror.AnErrWeakPassword) {
		t.Errorf("Esperado erro %v, obtido %v", msgerror.AnErrWeakPassword, err)
	}
}

func TestUser_WithName_Success(t *testing.T) {
	user, _ := entity.CreateUser(
		"Original",
		"original@example.com",
		"Pass123!",
		"",
	)

	newName, _ := vo.NewName("New Name", 3, 50)
	updatedUser, err := user.WithName(newName)
	if err != nil {
		t.Fatalf("WithName retornou erro inesperado: %v", err)
	}

	if updatedUser.Name != newName {
		t.Errorf("Nome não atualizado. Esperado: %v, Obtido: %v", newName, updatedUser.Name)
	}
}

func TestUser_WithPasswordHash_Success(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"OriginalPass123!",
		"",
	)

	newHash, _ := vo.NewPasswordHash("NewPass123!")
	updatedUser, err := user.WithPasswordHash(newHash)
	if err != nil {
		t.Fatalf("WithPasswordHash retornou erro inesperado: %v", err)
	}

	if updatedUser.PasswordHash != newHash {
		t.Errorf("Hash não atualizado. Esperado: %v, Obtido: %v", newHash, updatedUser.PasswordHash)
	}
}

func TestUser_PasswordResetToken(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)

	if err := user.GeneratePasswordResetToken(); err != nil {
		t.Fatalf("GeneratePasswordResetToken falhou: %v", err)
	}

	if user.PasswordResetToken == "" {
		t.Error("Token não foi gerado")
	}

	if time.Until(user.PasswordResetExpires).Round(time.Minute) < 55*time.Minute {
		t.Errorf("Tempo de expiração inválido: %v", user.PasswordResetExpires)
	}

	user.ClearResetToken()
	if user.PasswordResetToken != "" {
		t.Error("Token não foi limpo")
	}
	if !user.PasswordResetExpires.IsZero() {
		t.Error("Expiração não foi resetada")
	}
}

func TestCreateUser_WithImageURL(t *testing.T) {
	user, err := entity.CreateUser(
		"John",
		"john@example.com",
		"Pass123!",
		"https://example.com/avatar.jpg",
	)
	if err != nil {
		t.Fatalf("CreateUser falhou: %v", err)
	}

	expectedURL := "https://example.com/avatar.jpg"
	if user.ImageURL.String() != expectedURL {
		t.Errorf("URL de imagem incorreta: esperado '%s', obtido '%s'",
			expectedURL, user.ImageURL.String())
	}
}

func TestNewUser_Valid(t *testing.T) {
	id := vo.NewID()
	name, _ := vo.NewName("Test", 3, 50)
	email, _ := vo.NewEmail("test@example.com")
	hash, _ := vo.NewPasswordHash("Pass123!")
	url, _ := vo.NewURL("https://example.com/image.jpg")

	user, err := entity.NewUser(id, name, email, hash, &url)
	if err != nil {
		t.Fatalf("NewUser falhou: %v", err)
	}

	if user.ID != id {
		t.Errorf("ID incorreto: esperado %v, obtido %v", id, user.ID)
	}
}

func TestCreateUser_WithoutImageURL(t *testing.T) {
	user, err := entity.CreateUser(
		"John",
		"john@example.com",
		"Pass123!",
		"",
	)
	if err != nil {
		t.Fatalf("CreateUser falhou: %v", err)
	}

	if !user.ImageURL.IsEmpty() {
		t.Errorf("URL de imagem deveria ser vazia, mas é '%s'", user.ImageURL.String())
	}
}

func TestUser_GeneratePasswordResetToken_Error(t *testing.T) {
	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	rand.Reader = errorReader{}

	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)

	err := user.GeneratePasswordResetToken()
	if err == nil {
		t.Fatal("Esperava um erro, mas não ocorreu")
	}

	expectedErr := "erro simulado na leitura"
	if err.Error() != expectedErr {
		t.Errorf("Erro incorreto: esperado '%s', obtido '%s'", expectedErr, err.Error())
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("erro simulado na leitura")
}
