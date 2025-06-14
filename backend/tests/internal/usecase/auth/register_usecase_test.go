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

	assert.ErrorContains(t, err, "Invalid Password and confirmation")
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

	assert.ErrorContains(t, err, "invalid name")
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

	assert.ErrorContains(t, err, "invalid email")
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
		Password:             "any",
		PasswordConfirmation: "any", // Adicionado confirmação
	})

	assert.ErrorIs(t, err, msgerror.AnErrUserExists)
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
		{"Too Short", "ab", "invalid name"},
		{"Too Long", strings.Repeat("a", 101), "invalid name"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
			err := handler.Execute(context.Background(), dto.RegisterParams{
				Name:  tc.inputName,
				Email: "test@test.com",
			})

			assert.ErrorContains(t, err, tc.expectedError)
			mockRepo.AssertNotCalled(t, "GetByEmail")
			mockCrypto.AssertNotCalled(t, "Encrypt")
		})
	}
}

func TestRegisterWithInvalidPasswordHash(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	// Hash que não é bcrypt e tem menos de 8 caracteres
	invalidHash := "short"
	mockCrypto.On("Encrypt", "any").Return(invalidHash, nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "any",
		PasswordConfirmation: "any",
	})

	assert.ErrorContains(t, err, "failed to secure password")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestRegisterSuccessfully(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	name, _ := vo.NewName("New User", 0, 0)
	passwordHash, _ := vo.NewPasswordHash("hashed-password")
	email, _ := vo.NewEmail("new@test.com")

	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	newUser := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		ImageURL:     vo.URL{},
		CreatedAt:    time.Now(),
	}
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(newUser, nil)
	mockCrypto.On("Encrypt", "valid-password").Return("hashed-password", nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "New User",
		Email:                "new@test.com",
		Password:             "valid-password",
		PasswordConfirmation: "valid-password",
	})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
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
		Password:             "any",
		PasswordConfirmation: "any", // Adicionado confirmação
	})

	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithEncryptError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("", errors.New("encryption failed"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "any",
		PasswordConfirmation: "any", // Adicionado confirmação
	})

	assert.ErrorContains(t, err, "failed to secure password")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithPasswordHashError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("", nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "any",
		PasswordConfirmation: "any", // Adicionado confirmação
	})

	assert.ErrorContains(t, err, "failed to secure password")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithSaveError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("hashed-password", nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "any",
		PasswordConfirmation: "any", // Adicionado confirmação
	})

	assert.ErrorContains(t, err, "failed to create user")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWhenSaveReturnsNilUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("hashed-password", nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil, nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "test@test.com",
		Password:             "any",
		PasswordConfirmation: "any", // Adicionado confirmação
	})

	assert.ErrorIs(t, err, msgerror.AnErrNoSavedUser)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithCustomNameLimits(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:  "AB",
		Email: "test@test.com",
	})

	assert.ErrorContains(t, err, "invalid name")
}

func TestRegisterWithInvalidImageURL(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	// Configurar o mock para GetByEmail
	email, _ := vo.NewEmail("valid@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:                 "Valid Name",
		Email:                "valid@test.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
		ImageURL:             "invalid-url", // URL inválida
	})

	assert.ErrorContains(t, err, "invalid image URL")
	mockRepo.AssertExpectations(t) // Verifica que GetByEmail foi chamado
	mockCrypto.AssertNotCalled(t, "Encrypt")
	mockRepo.AssertNotCalled(t, "Save")
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
		ImageURL:             "https://example.com/image.jpg", // URL válida
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
		ImageURL:             "", // URL vazia
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}
