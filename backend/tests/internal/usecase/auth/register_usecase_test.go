package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterWithPasswordMismatch(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "valid@test.com",
		Password:             "password",
		PasswordConfirmation: "different",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "passwords do not match", valErr.FieldErrors["password_confirmation"])

	// Não deve chamar GetByEmail porque há erro de validação
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithShortPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "valid@test.com",
		Password:             "short",
		PasswordConfirmation: "short",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "must be at least 8 characters", valErr.FieldErrors["password"])
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithInvalidName(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name: "A",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Contains(t, valErr.FieldErrors["name"], "name too short")
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithInvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:  "Valid Name",
		Email: "invalid-email",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Contains(t, valErr.FieldErrors["email"], "invalid email format")
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithExistingUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("existing@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return(&entity.User{}, nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Existing User",
		Email:                "existing@test.com",
		Password:             "valid-password", // Usar senha válida
		PasswordConfirmation: "valid-password", // Igual à senha
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, msgerror.AnErrUserExists.Error(), valErr.FieldErrors["email"])
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithEdgeCaseNameError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	testCases := []struct {
		name          string
		inputName     string
		expectedError string
	}{
		{"Too Short", "ab", "name too short"},
		{"Too Long", strings.Repeat("a", 101), "name too long"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
			err := handler.Execute(context.Background(), dto.RegisterParams{
				Name:  tc.inputName,
				Email: "test@test.com",
			})

			var valErr *msgerror.ValidationErrors
			assert.ErrorAs(t, err, &valErr)
			assert.Contains(t, valErr.FieldErrors["name"], tc.expectedError)
			mockRepo.AssertNotCalled(t, "GetByEmail")
			mockCrypto.AssertNotCalled(t, "Encrypt")
		})
	}
}

func TestRegisterWithUnexpectedErrorOnGetByEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	expectedErr := errors.New("unexpected error")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), expectedErr)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "valid-password", // Usar senha válida
		PasswordConfirmation: "valid-password", // Igual à senha
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check email existence")
	assert.Contains(t, err.Error(), "unexpected error")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithEncryptError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	mockCrypto.On("Encrypt", "valid-password").Return("", errors.New("encryption failed"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to secure password")
	assert.Contains(t, err.Error(), "encryption failed")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithPasswordHashError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	// Forçar erro na criação do PasswordHash
	mockCrypto.On("Encrypt", "valid-password").Return("", errors.New("invalid hash format"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to secure password") // Atualizar mensagem esperada
	assert.Contains(t, err.Error(), "invalid hash format")

	// Garantir que o Save não foi chamado
	mockRepo.AssertNotCalled(t, "Save")
}

func TestRegisterWithEmptyPasswordHash(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	// Simular situação onde o Encrypt retorna string vazia sem erro
	mockCrypto.On("Encrypt", "valid-password").Return("", nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create password hash")
	assert.Contains(t, err.Error(), "password must be at least 8 characters")

	// Garantir que o Save não foi chamado
	mockRepo.AssertNotCalled(t, "Save")
}

func TestRegisterWithSaveError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	mockCrypto.On("Encrypt", "valid-password").Return("hashed-password", nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWhenSaveReturnsNilUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	// DADOS VÁLIDOS PARA PASSAR NAS VALIDAÇÕES INICIAIS
	validName := "Valid Name"
	validEmail := "test@test.com"
	validPassword := "valid-password123!"

	// Configurar mocks para fluxo completo
	email, _ := vo.NewEmail(validEmail)
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", validPassword).Return("hashed-password", nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil, nil) // Simular retorno nil do Save

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 validName,
		Email:                validEmail,
		Password:             validPassword,
		PasswordConfirmation: validPassword,
	})

	// Verificar erro específico
	assert.ErrorIs(t, err, msgerror.AnErrNoSavedUser)

	// Verificar chamadas dos mocks
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithCustomNameLimits(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:  "AB", // Muito curto
		Email: "test@test.com",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Contains(t, valErr.FieldErrors["name"], "name too short")
}

func TestRegisterWithInvalidImageURL(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "valid@test.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
		ImageURL:             "invalid-url",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Contains(t, valErr.FieldErrors["image_url"], "invalid URL format")

	// Não deve chamar operações posteriores
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithValidImageURL(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	// Configurar os mocks necessários
	name, _ := vo.NewName("New User", 0, 0)
	passwordHash, _ := vo.NewPasswordHash("hashed-password")
	email, _ := vo.NewEmail("new@test.com")
	validURL, _ := vo.NewURL("https://example.com/image.jpg")

	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "valid-password").Return("hashed-password", nil)

	newUser := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		ImageURL:     validURL,
		CreatedAt:    time.Now(),
	}
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(newUser, nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "New User",
		Email:                "new@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
		ImageURL:             "https://example.com/image.jpg",
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithEmptyImageURL(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	// Configurar os mocks necessários
	name, _ := vo.NewName("New User", 0, 0)
	passwordHash, _ := vo.NewPasswordHash("hashed-password")
	email, _ := vo.NewEmail("new@test.com")

	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "valid-password").Return("hashed-password", nil)

	newUser := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		ImageURL:     vo.URL{},
		CreatedAt:    time.Now(),
	}
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(newUser, nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "New User",
		Email:                "new@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
		ImageURL:             "",
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithEmptyPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "valid@test.com",
		Password:             "",
		PasswordConfirmation: "",
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "must be at least 8 characters", valErr.FieldErrors["password"])
}

func TestRegisterWithMultipleValidationErrors(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "A",         // Nome muito curto
		Email:                "invalid",   // E-mail inválido
		Password:             "short",     // Senha curta
		PasswordConfirmation: "different", // Confirmação diferente
	})

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Len(t, valErr.FieldErrors, 4)
	assert.Contains(t, valErr.FieldErrors["name"], "name too short")
	assert.Contains(t, valErr.FieldErrors["email"], "invalid email format")
	assert.Equal(t, "must be at least 8 characters", valErr.FieldErrors["password"])
	assert.Equal(t, "passwords do not match", valErr.FieldErrors["password_confirmation"])
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithValidDataButCryptoError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "valid-password").Return("", errors.New("crypto error"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to secure password")
	assert.Contains(t, err.Error(), "crypto error")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}
