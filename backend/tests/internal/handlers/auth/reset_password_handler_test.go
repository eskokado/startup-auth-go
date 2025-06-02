package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResetPasswordHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode) // Adiciona para suprimir logs

	tests := []struct {
		name         string
		input        interface{}
		mockError    error
		expectedCode int
		expectedBody gin.H
	}{
		{
			name:         "invalid request body",
			input:        "{invalid}",
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid request body"},
		},
		{
			name:         "invalid token",
			input:        dto.ResetPasswordInput{Token: "invalid", Password: "newpass"},
			mockError:    msgerror.AnErrInvalidToken,
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": msgerror.AnErrInvalidToken.Error()},
		},
		{
			name:         "expired token",
			input:        dto.ResetPasswordInput{Token: "expired", Password: "newpass"},
			mockError:    msgerror.AnErrExpiredToken,
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": msgerror.AnErrExpiredToken.Error()},
		},
		{
			name:         "internal server error",
			input:        dto.ResetPasswordInput{Token: "valid", Password: "newpass"},
			mockError:    assert.AnError,
			expectedCode: http.StatusInternalServerError,
			expectedBody: gin.H{"error": "failed to reset password"},
		},
		{
			name:         "success",
			input:        dto.ResetPasswordInput{Token: "valid", Password: "newpass"},
			mockError:    nil,
			expectedCode: http.StatusNoContent, // Corrigido para 204
			expectedBody: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configura o mock corretamente
			uc := &mocks.MockResetPasswordUseCase{
				Err: tt.mockError, // Injeta o erro
			}

			handler := handlers.NewResetPasswordHandler(uc)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Cria a requisição corretamente mesmo para corpo inválido
			var req *http.Request
			switch v := tt.input.(type) {
			case string:
				req = httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader([]byte(v)))
			default:
				body, _ := json.Marshal(tt.input)
				req = httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.Handle(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				var response gin.H
				_ = json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedBody, response)
			} else if w.Code != http.StatusNoContent {
				assert.Empty(t, w.Body.String(), "Response body should be empty")
			}
		})
	}
}
