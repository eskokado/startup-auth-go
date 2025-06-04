package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso - Login válido", func(t *testing.T) {
		mockUseCase := new(mocks.MockLoginUseCase)
		handler := handlers.NewLoginHandler(mockUseCase)

		fixedID, _ := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		email, _ := vo.NewEmail("test@example.com")

		expectedLoginResult := dto.LoginResult{
			UserID: fixedID,
			Email:  email,
			Token:  "token_jwt_gerado_pelo_use_case",
		}

		mockUseCase.On("Execute", mock.Anything, "test@example.com", "senha123").
			Return(expectedLoginResult, nil)

		reqBody := `{"email": "test@example.com", "password": "senha123"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/login", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"access_token": "token_jwt_gerado_pelo_use_case"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Erro - Body inválido", func(t *testing.T) {
		handler := handlers.NewLoginHandler(nil)

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

		handler := handlers.NewLoginHandler(mockUseCase)

		reqBody := `{"email": "invalid@example.com", "password": "wrong"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/login", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
