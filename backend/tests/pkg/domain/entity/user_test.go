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

// --- Helper para simular erros no crypto/rand ---
type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("erro simulado na leitura")
}

// --- Testes de criação de usuário ---
func TestCreateUser_Valid(t *testing.T) {
	_, err := entity.CreateUser(
		"John Doe",
		"john@example.com",
		"SecurePass123!",
		"https://example.com/avatar.jpg",
	)
	if err != nil {
		t.Errorf("CreateUser() falhou: %v", err)
	}
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	_, err := entity.CreateUser(
		"John Doe",
		"invalid-email",
		"SecurePass123!",
		"",
	)

	var valErr *msgerror.ValidationErrors
	if !errors.As(err, &valErr) {
		t.Fatalf("Esperado ValidationErrors, recebido: %T", err)
	}

	if len(valErr.FieldErrors) == 0 {
		t.Error("Sem erros de campo no ValidationErrors")
	}
}

func TestCreateUser_WeakPassword(t *testing.T) {
	_, err := entity.CreateUser(
		"John Doe",
		"john@example.com",
		"short",
		"",
	)

	var valErr *msgerror.ValidationErrors
	if !errors.As(err, &valErr) {
		t.Fatalf("Esperado ValidationErrors, recebido: %T", err)
	}

	if len(valErr.FieldErrors) == 0 {
		t.Error("Sem erros de campo no ValidationErrors")
	}
}

func TestCreateUser_InvalidImageURL(t *testing.T) {
	_, err := entity.CreateUser(
		"John Doe",
		"john@example.com",
		"SecurePass123!",
		"htp://invalid-url",
	)

	var valErr *msgerror.ValidationErrors
	if !errors.As(err, &valErr) {
		t.Fatalf("Esperado ValidationErrors, recebido: %T", err)
	}
}

// --- Testes de igualdade entre usuários ---
func TestUser_Equal_SameUser(t *testing.T) {
	user, _ := entity.CreateUser(
		"Alice",
		"alice@example.com",
		"StrongPass123!",
		"",
	)

	if !user.Equal(user) {
		t.Error("Equal() deve retornar true para o mesmo usuário")
	}
}

func TestUser_Equal_DifferentUsers(t *testing.T) {
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

	if user1.Equal(user2) {
		t.Error("Equal() deve retornar false para usuários diferentes")
	}
}

// --- Testes de verificação de senha ---
func TestUser_VerifyPassword_Correct(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test User",
		"test@example.com",
		"CorrectPass123!",
		"",
	)

	if !user.VerifyPassword("CorrectPass123!") {
		t.Error("VerifyPassword() deve retornar true para senha correta")
	}
}

func TestUser_VerifyPassword_Incorrect(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test User",
		"test@example.com",
		"CorrectPass123!",
		"",
	)

	if user.VerifyPassword("WrongPass456@") {
		t.Error("VerifyPassword() deve retornar false para senha incorreta")
	}
}

// --- Testes de validação no construtor ---
func TestNewUser_EmptyEmail(t *testing.T) {
	validID := vo.NewID()
	validName, _ := vo.NewName("Test", 3, 50)
	emptyEmail, _ := vo.NewEmail("")
	validHash, _ := vo.NewPasswordHash("ValidPass123!")

	_, err := entity.NewUser(
		validID,
		validName,
		emptyEmail,
		validHash,
		nil,
	)

	if !errors.Is(err, msgerror.AnErrEmptyEmail) {
		t.Errorf("Esperado %v, recebido %v", msgerror.AnErrEmptyEmail, err)
	}
}

func TestNewUser_EmptyPasswordHash(t *testing.T) {
	validID := vo.NewID()
	validName, _ := vo.NewName("Test", 3, 50)
	validEmail, _ := vo.NewEmail("test@example.com")
	emptyHash, _ := vo.NewPasswordHash("")

	_, err := entity.NewUser(
		validID,
		validName,
		validEmail,
		emptyHash,
		nil,
	)

	if !errors.Is(err, msgerror.AnErrWeakPassword) {
		t.Errorf("Esperado %v, recebido %v", msgerror.AnErrWeakPassword, err)
	}
}

// --- Testes de atualização de nome ---
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
		t.Fatalf("WithName falhou: %v", err)
	}

	if updatedUser.Name != newName {
		t.Errorf("Nome não atualizado: esperado %v, recebido %v", newName, updatedUser.Name)
	}
}

func TestUser_WithName_EmptyName(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)

	emptyName, _ := vo.NewName("", 0, 50)
	_, err := user.WithName(emptyName)

	if !errors.Is(err, msgerror.AnErrInvalidName) {
		t.Errorf("Esperado %v, recebido %v", msgerror.AnErrInvalidName, err)
	}
}

func TestUser_WithName_SameName(t *testing.T) {
	user, _ := entity.CreateUser(
		"Original",
		"original@example.com",
		"Pass123!",
		"",
	)

	sameName, _ := vo.NewName("Original", 3, 50)
	_, err := user.WithName(sameName)

	if !errors.Is(err, msgerror.AnErrNameDifferent) {
		t.Errorf("Esperado %v, recebido %v", msgerror.AnErrNameDifferent, err)
	}
}

// --- Testes de atualização de senha ---
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
		t.Fatalf("WithPasswordHash falhou: %v", err)
	}

	if updatedUser.PasswordHash != newHash {
		t.Error("Hash da senha não atualizado")
	}
}

func TestUser_WithPasswordHash_EmptyHash(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)

	emptyHash, _ := vo.NewPasswordHash("")
	_, err := user.WithPasswordHash(emptyHash)

	if !errors.Is(err, msgerror.AnErrWeakPassword) {
		t.Errorf("Esperado %v, recebido %v", msgerror.AnErrWeakPassword, err)
	}
}

// --- Testes de token de reset de senha ---
func TestUser_PasswordResetToken_GenerateAndClear(t *testing.T) {
	user, _ := entity.CreateUser(
		"Test",
		"test@example.com",
		"Pass123!",
		"",
	)

	// Geração do token
	if err := user.GeneratePasswordResetToken(); err != nil {
		t.Fatalf("GeneratePasswordResetToken falhou: %v", err)
	}

	if user.PasswordResetToken == "" {
		t.Error("Token não gerado")
	}

	if time.Until(user.PasswordResetExpires) < 55*time.Minute {
		t.Error("Tempo de expiração inválido")
	}

	// Limpeza do token
	user.ClearResetToken()
	if user.PasswordResetToken != "" {
		t.Error("Token não limpo")
	}
	if !user.PasswordResetExpires.IsZero() {
		t.Error("Expiração não resetada")
	}
}

func TestUser_GeneratePasswordResetToken_Error(t *testing.T) {
	// Salvar e restaurar o leitor original
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
		t.Fatal("Esperado erro, não ocorreu")
	}

	expectedErr := "erro simulado na leitura"
	if err.Error() != expectedErr {
		t.Errorf("Esperado %v, recebido %v", expectedErr, err)
	}
}

// --- Testes de URL de imagem ---
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
		t.Errorf("URL incorreta: esperado %s, recebido %s", expectedURL, user.ImageURL.String())
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
		t.Errorf("URL deveria ser vazia: %s", user.ImageURL.String())
	}
}

// --- Teste de construtor completo ---
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

	if user.ID != id || user.Name != name || user.Email != email || user.PasswordHash != hash || user.ImageURL != url {
		t.Error("Propriedades do usuário incorretas")
	}
}
