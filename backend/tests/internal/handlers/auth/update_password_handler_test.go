package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUpdatePasswordUseCase struct {
	mock.Mock
}

func (m *MockUpdatePasswordUseCase) Execute(ctx context.Context, userID vo.ID, currentPassword, newPassword string) error {
	args := m.Called(ctx, userID, currentPassword, newPassword)
	return args.Error(0)
}

func TestUpdatePasswordHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Password updated", func(t *testing.T) {
		mockUseCase := new(MockUpdatePasswordUseCase)
		handler := handlers.NewUpdatePasswordHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		mockUseCase.On("Execute", mock.Anything, userID, "oldPass", "newPass").Return(nil)

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"current_password": "oldPass", "new_password": "newPass"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Invalid user ID in context", func(t *testing.T) {
		handler := handlers.NewUpdatePasswordHandler(nil)

		router := gin.Default()
		router.PUT("/password", handler.Handle)

		reqBody := `{"current_password": "old", "new_password": "new"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		handler := handlers.NewUpdatePasswordHandler(nil)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(`{invalid`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Error - Current password mismatch", func(t *testing.T) {
		mockUseCase := new(MockUpdatePasswordUseCase)
		handler := handlers.NewUpdatePasswordHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		mockUseCase.On("Execute", mock.Anything, userID, "wrong", "newPass").Return(msgerror.AnErrInvalidCredentials)

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"current_password": "wrong", "new_password": "newPass"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error": "invalid credentials"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Weak new password", func(t *testing.T) {
		mockUseCase := new(MockUpdatePasswordUseCase)
		handler := handlers.NewUpdatePasswordHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		mockUseCase.On("Execute", mock.Anything, userID, "oldPass", "weak").Return(msgerror.AnErrWeakPassword)

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"current_password": "oldPass", "new_password": "weak"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "password does not meet security requirements"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Invalid user ID format in context", func(t *testing.T) {
		handler := handlers.NewUpdatePasswordHandler(nil)

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", "invalid-uuid-format") // ID inválido
			handler.Handle(c)
		})

		reqBody := `{"current_password": "old", "new_password": "new"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "invalid ID format"}`, resp.Body.String())
	})

	t.Run("Error - User not found", func(t *testing.T) {
		mockUseCase := new(MockUpdatePasswordUseCase)
		handler := handlers.NewUpdatePasswordHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		// Simula erro de usuário não encontrado
		mockUseCase.On("Execute", mock.Anything, userID, "oldPass", "newPass").Return(msgerror.AnErrUserNotFound)

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"current_password": "oldPass", "new_password": "newPass"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error": "user not found"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Unhandled internal error", func(t *testing.T) {
		mockUseCase := new(MockUpdatePasswordUseCase)
		handler := handlers.NewUpdatePasswordHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		// Simula erro não mapeado
		mockUseCase.On("Execute", mock.Anything, userID, "oldPass", "newPass").Return(assert.AnError)

		router := gin.Default()
		router.PUT("/password", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"current_password": "oldPass", "new_password": "newPass"}`
		req, _ := http.NewRequest(http.MethodPut, "/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error": "failed to update password"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})
}
