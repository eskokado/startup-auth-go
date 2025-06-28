package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestForgotPasswordHandler_InvalidRequestBody(t *testing.T) {
	// Setup
	uc := new(mocks.MockForgotPasswordUseCase)
	handler := handlers.NewForgotPasswordHandler(uc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request com corpo inválido
	body := []byte("{invalid}")
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	handler.Handle(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"invalid character 'i' looking for beginning of object key string"}`, w.Body.String())
}

func TestForgotPasswordHandler_InvalidEmailFormat(t *testing.T) {
	// Setup
	uc := new(mocks.MockForgotPasswordUseCase)
	handler := handlers.NewForgotPasswordHandler(uc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request com email inválido
	input := dto.ForgotPasswordInput{Email: "invalid"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	handler.Handle(c)

	// Assert - Mantém a mensagem para formato inválido
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"invalid email format"}`, w.Body.String())
}

func TestForgotPasswordHandler_InternalServerError(t *testing.T) {
	// Setup
	uc := new(mocks.MockForgotPasswordUseCase)
	handler := handlers.NewForgotPasswordHandler(uc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Configura o mock para retornar erro
	uc.On("Execute", mock.Anything, mock.Anything).Return(assert.AnError)

	// Request válida
	input := dto.ForgotPasswordInput{Email: "valid@example.com"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	handler.Handle(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"error":"assert.AnError general error for testing"}`, w.Body.String())
	uc.AssertExpectations(t)
}

func TestForgotPasswordHandler_Success(t *testing.T) {
	// Setup
	uc := new(mocks.MockForgotPasswordUseCase)
	handler := handlers.NewForgotPasswordHandler(uc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Configura o mock para sucesso
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil)

	// Request válida
	input := dto.ForgotPasswordInput{Email: "valid@example.com"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	handler.Handle(c)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	uc.AssertExpectations(t)
}

func TestForgotPasswordHandler_EmptyEmail(t *testing.T) {
	// Setup
	uc := new(mocks.MockForgotPasswordUseCase)
	handler := handlers.NewForgotPasswordHandler(uc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request com email vazio
	input := dto.ForgotPasswordInput{Email: ""}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	handler.Handle(c)

	// Assert - Agora espera a mensagem específica para email vazio
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"email cannot be empty"}`, w.Body.String())
}
func TestForgotPasswordHandler_EmptyBody(t *testing.T) {
	// Setup
	uc := new(mocks.MockForgotPasswordUseCase)
	handler := handlers.NewForgotPasswordHandler(uc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request sem corpo
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", nil)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	handler.Handle(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"EOF"}`, w.Body.String())
}
