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
	"github.com/stretchr/testify/mock"
)

func TestResetPasswordHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		input        interface{}
		mockSetup    func(uc *mocks.MockResetPasswordUseCase)
		expectedCode int
		expectedBody string
	}{
		{
			name:  "invalid request body",
			input: "{invalid}",
			mockSetup: func(uc *mocks.MockResetPasswordUseCase) {
				// Não configurar mock pois não deve ser chamado
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid request body"}`,
		},
		{
			name:  "invalid token",
			input: dto.ResetPasswordInput{Token: "invalid", Password: "newpass"},
			mockSetup: func(uc *mocks.MockResetPasswordUseCase) {
				uc.On("Execute", mock.Anything, "invalid", "newpass").
					Return(msgerror.AnErrInvalidToken)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid token"}`,
		},
		{
			name:  "expired token",
			input: dto.ResetPasswordInput{Token: "expired", Password: "newpass"},
			mockSetup: func(uc *mocks.MockResetPasswordUseCase) {
				uc.On("Execute", mock.Anything, "expired", "newpass").
					Return(msgerror.AnErrExpiredToken)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"expired token"}`,
		},
		{
			name:  "internal server error",
			input: dto.ResetPasswordInput{Token: "valid", Password: "newpass"},
			mockSetup: func(uc *mocks.MockResetPasswordUseCase) {
				uc.On("Execute", mock.Anything, "valid", "newpass").
					Return(assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"failed to reset password"}`,
		},
		{
			name:  "success",
			input: dto.ResetPasswordInput{Token: "valid", Password: "newpass"},
			mockSetup: func(uc *mocks.MockResetPasswordUseCase) {
				uc.On("Execute", mock.Anything, "valid", "newpass").
					Return(nil)
			},
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cria mock do caso de uso
			uc := new(mocks.MockResetPasswordUseCase)
			if tt.mockSetup != nil {
				tt.mockSetup(uc)
			}

			// Cria handler com o mock
			handler := handlers.NewResetPasswordHandler(uc)

			// Configura o roteador Gin
			router := gin.New()
			router.POST("/reset-password", handler.Handle)

			// Cria a requisição
			var bodyBytes []byte
			switch v := tt.input.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Executa a requisição
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificações
			assert.Equal(t, tt.expectedCode, w.Code, "Código de status HTTP incorreto")
			assert.Equal(t, tt.expectedBody, w.Body.String(), "Corpo da resposta incorreto")

			// Verifica se o mock foi chamado conforme esperado
			uc.AssertExpectations(t)
		})
	}
}
