package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso - Login válido", func(t *testing.T) {
		// Configurar mocks
		mockUseCase := new(mocks.MockLoginUseCase)
		mockToken := new(mocks.MockTokenProvider)

		handler := handlers.NewLoginHandler(mockUseCase, mockToken)

		fixedID, _ := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		email, _ := vo.NewEmail("test@example.com")

		// Dados de teste
		expectedLoginResult := dto.LoginResult{
			UserID: fixedID,
			Email:  email,
		}
		expectedClaims := providers.Claims{
			UserID: fixedID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test@example.com",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}
		mockUseCase.On("Execute", mock.Anything, "test@example.com", "senha123").
			Return(expectedLoginResult, nil)
		mockToken.On("Generate", expectedClaims).
			Return("token_jwt", nil)

		// Configurar requisição
		reqBody := `{"email": "test@example.com", "password": "senha123"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Executar
		router := gin.Default()
		router.POST("/login", handler.Handle)
		router.ServeHTTP(resp, req)

		// Validar
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"access_token": "token_jwt"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("Erro - Body inválido", func(t *testing.T) {
		handler := handlers.NewLoginHandler(nil, nil)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{invalid`))
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/login", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Erro - Credenciais inválidas", func(t *testing.T) {
		mockUseCase := new(mocks.MockLoginUseCase)
		mockUseCase.On("Execute", mock.Anything, "invalid@example.com", "wrong").
			Return(dto.LoginResult{}, errors.New("credenciais inválidas"))

		handler := handlers.NewLoginHandler(mockUseCase, nil)

		reqBody := `{"email": "invalid@example.com", "password": "wrong"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/login", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Erro - Falha ao gerar token", func(t *testing.T) {
		// Configurar mocks
		mockUseCase := new(mocks.MockLoginUseCase)
		mockToken := new(mocks.MockTokenProvider)

		handler := handlers.NewLoginHandler(mockUseCase, mockToken)

		fixedID, _ := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		email, _ := vo.NewEmail("test@example.com")

		// Configurar retorno válido do caso de uso
		expectedLoginResult := dto.LoginResult{
			UserID: fixedID,
			Email:  email,
		}
		mockUseCase.On("Execute", mock.Anything, "test@example.com", "senha123").
			Return(expectedLoginResult, nil)

		// Configurar token provider para retornar erro
		expectedClaims := providers.Claims{
			UserID: fixedID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test@example.com",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}
		mockToken.On("Generate", expectedClaims).
			Return("", errors.New("erro na geração do token"))

		// Configurar requisição
		reqBody := `{"email": "test@example.com", "password": "senha123"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Executar
		router := gin.Default()
		router.POST("/login", handler.Handle)
		router.ServeHTTP(resp, req)

		// Validar
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"failed to generate token"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})
}
