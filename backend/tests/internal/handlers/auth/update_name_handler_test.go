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

type MockUpdateNameUseCase struct {
	mock.Mock
}

func (m *MockUpdateNameUseCase) Execute(ctx context.Context, userID vo.ID, newName string) error {
	args := m.Called(ctx, userID, newName)
	return args.Error(0)
}

func TestUpdateNameHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Name updated", func(t *testing.T) {
		mockUseCase := new(MockUpdateNameUseCase)
		handler := handlers.NewUpdateNameHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		mockUseCase.On("Execute", mock.Anything, userID, "New Name").Return(nil)

		router := gin.Default()
		router.PUT("/name", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"name": "New Name"}`
		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Invalid user ID in context", func(t *testing.T) {
		handler := handlers.NewUpdateNameHandler(nil)

		router := gin.Default()
		router.PUT("/name", handler.Handle)

		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(`{"name": "New Name"}`))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		handler := handlers.NewUpdateNameHandler(nil)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		router := gin.Default()
		router.PUT("/name", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(`{invalid`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Error - Use case returns invalid name", func(t *testing.T) {
		mockUseCase := new(MockUpdateNameUseCase)
		handler := handlers.NewUpdateNameHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		mockUseCase.On("Execute", mock.Anything, userID, "Short").Return(msgerror.AnErrNameTooShort)

		router := gin.Default()
		router.PUT("/name", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"name": "Short"}`
		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "name too short"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Unhandled internal error", func(t *testing.T) {
		mockUseCase := new(MockUpdateNameUseCase)
		handler := handlers.NewUpdateNameHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		// Simula erro não mapeado
		mockUseCase.On("Execute", mock.Anything, userID, "newName").Return(assert.AnError)

		router := gin.Default()
		router.PUT("/name", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"name": "newName"}`
		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error": "failed to update name"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		mockUseCase := new(MockUpdateNameUseCase)
		handler := handlers.NewUpdateNameHandler(mockUseCase)

		userID, err := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			t.Fatal(err)
		}

		// Simula erro de usuário não encontrado
		mockUseCase.On("Execute", mock.Anything, userID, "newName").Return(msgerror.AnErrUserNotFound)

		router := gin.Default()
		router.PUT("/name", func(c *gin.Context) {
			c.Set("userID", userID.String())
			handler.Handle(c)
		})

		reqBody := `{"name": "newName"}`
		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error": "user not found"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error - Invalid user ID format in context", func(t *testing.T) {
		handler := handlers.NewUpdateNameHandler(nil)

		router := gin.Default()
		router.PUT("/name", func(c *gin.Context) {
			c.Set("userID", "invalid-uuid-format") // ID inválido
			handler.Handle(c)
		})

		reqBody := `{"name": "New Name"}`
		req, _ := http.NewRequest(http.MethodPut, "/name", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "invalid ID format"}`, resp.Body.String())
	})
}
