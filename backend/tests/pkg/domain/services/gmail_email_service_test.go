package service_test

import (
	"errors"
	"os"
	"testing"

	service "github.com/eskokado/startup-auth-go/backend/pkg/domain/services"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmailService_SendResetPasswordEmail(t *testing.T) {
	// Configuração inicial para garantir ambiente limpo
	originalEnv := map[string]string{
		"FROM_EMAIL":         os.Getenv("FROM_EMAIL"),
		"FRONTEND_RESET_URL": os.Getenv("FRONTEND_RESET_URL"),
		"SMTP_HOST":          os.Getenv("SMTP_HOST"),
		"SMTP_PORT":          os.Getenv("SMTP_PORT"),
		"SMTP_USERNAME":      os.Getenv("SMTP_USERNAME"),
		"SMTP_PASSWORD":      os.Getenv("SMTP_PASSWORD"),
	}

	defer func() {
		for k, v := range originalEnv {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// Testes individuais com ambiente controlado
	t.Run("Sucesso - Email enviado", func(t *testing.T) {
		os.Setenv("FROM_EMAIL", "no-reply@example.com")
		os.Setenv("FRONTEND_RESET_URL", "https://app.example.com/reset-password")

		mockSender := new(mocks.MockSenderService)
		emailService := service.NewEmailService(mockSender)

		mockSender.On("DialAndSend", mock.Anything).Return(nil)

		email, _ := vo.NewEmail("user@example.com")
		err := emailService.SendResetPasswordEmail(email, "reset-token-123")

		assert.NoError(t, err)
		mockSender.AssertExpectations(t)
	})

	t.Run("Erro - Falha no envio", func(t *testing.T) {
		os.Setenv("FROM_EMAIL", "no-reply@example.com")
		os.Setenv("FRONTEND_RESET_URL", "https://app.example.com/reset-password")

		mockSender := new(mocks.MockSenderService)
		emailService := service.NewEmailService(mockSender)

		mockSender.On("DialAndSend", mock.Anything).Return(errors.New("smtp error"))

		email, _ := vo.NewEmail("user@example.com")
		err := emailService.SendResetPasswordEmail(email, "reset-token-123")

		assert.Error(t, err)
		assert.Equal(t, "smtp error", err.Error())
		mockSender.AssertExpectations(t)
	})

	t.Run("Erro - FROM_EMAIL não definido", func(t *testing.T) {
		os.Unsetenv("FROM_EMAIL")
		os.Setenv("FRONTEND_RESET_URL", "https://app.example.com/reset-password")

		mockSender := new(mocks.MockSenderService)
		emailService := service.NewEmailService(mockSender)

		email, _ := vo.NewEmail("user@example.com")
		err := emailService.SendResetPasswordEmail(email, "reset-token-123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FROM_EMAIL não está definido")
	})

	t.Run("Erro - FRONTEND_RESET_URL não definido", func(t *testing.T) {
		os.Setenv("FROM_EMAIL", "no-reply@example.com")
		os.Unsetenv("FRONTEND_RESET_URL")

		mockSender := new(mocks.MockSenderService)
		emailService := service.NewEmailService(mockSender)

		email, _ := vo.NewEmail("user@example.com")
		err := emailService.SendResetPasswordEmail(email, "reset-token-123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FRONTEND_RESET_URL não está definido")
	})

	t.Run("Sucesso - Cria dialer quando sender é nil", func(t *testing.T) {
		os.Setenv("SMTP_HOST", "smtp.example.com")
		os.Setenv("SMTP_PORT", "587")
		os.Setenv("SMTP_USERNAME", "user")
		os.Setenv("SMTP_PASSWORD", "pass")

		emailService := service.NewEmailService(nil)
		assert.NotNil(t, emailService)
	})
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		want    int
		wantErr bool
	}{
		{"Porta válida", "587", 587, false},
		{"Porta inválida", "invalid", 0, true},
		{"Porta vazia", "", 0, true},
		{"Porta zero", "0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ParsePort(tt.port)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewEmailService_Fallback(t *testing.T) {
	// Garantir ambiente limpo
	os.Unsetenv("FROM_EMAIL")
	os.Unsetenv("FRONTEND_RESET_URL")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "invalid") // Porta inválida
	os.Setenv("SMTP_USERNAME", "user")
	os.Setenv("SMTP_PASSWORD", "pass")

	t.Run("Usa fallback quando porta é inválida", func(t *testing.T) {
		emailService := service.NewEmailService(nil)
		assert.NotNil(t, emailService)
	})
}
