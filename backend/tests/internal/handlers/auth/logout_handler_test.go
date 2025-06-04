package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso - Logout", func(t *testing.T) {
		mockUseCase := new(mocks.MockLogoutUseCase)
		mockUseCase.On("Execute", mock.Anything, "valid_token").Return(nil)

		handler := handlers.NewLogoutHandler(mockUseCase)

		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/logout", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Erro - Sem header de autorização", func(t *testing.T) {
		handler := handlers.NewLogoutHandler(nil)

		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/logout", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Authorization header is required"}`, resp.Body.String())
	})

	t.Run("Erro - Formato de autorização inválido", func(t *testing.T) {
		handler := handlers.NewLogoutHandler(nil)

		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/logout", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Invalid authorization format"}`, resp.Body.String())
	})

	t.Run("Erro - Falha no caso de uso", func(t *testing.T) {
		mockUseCase := new(mocks.MockLogoutUseCase)
		mockUseCase.On("Execute", mock.Anything, "valid_token").Return(errors.New("erro"))

		handler := handlers.NewLogoutHandler(mockUseCase)

		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/logout", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"failed to logout"}`, resp.Body.String())
	})
}
