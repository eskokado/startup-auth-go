package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso - Registro válido", func(t *testing.T) {
		mockUseCase := new(mocks.MockRegisterUseCase)
		mockRepo := new(mocks.MockUserRepo)

		handler := handlers.NewRegisterHandler(mockUseCase, mockRepo)

		mockUseCase.On("Execute", mock.Anything, dto.RegisterParams{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "senha123",
		}).Return(nil)

		fixedID, _ := vo.ParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		fixedURL, _ := vo.NewURL("http://image.com")
		name, _ := vo.NewName("John Doe", 0, 0)
		email, _ := vo.NewEmail("john@example.com")
		mockRepo.On("GetByEmail", mock.Anything, email).
			Return(&entity.User{
				ID:       fixedID,
				Name:     name,
				Email:    email,
				ImageURL: fixedURL,
			}, nil)

		reqBody := `{
			"name": "John Doe",
			"email": "john@example.com",
			"password": "senha123",
			"image_url": "http://image.com"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/register", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.JSONEq(t, `{
			"id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"name": "John Doe",
			"email": "john@example.com",
			"image_url": "http://image.com",
			"password_reset_token": "",
			"password_reset_expires": ""
		}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro - Body inválido", func(t *testing.T) {
		handler := handlers.NewRegisterHandler(nil, nil)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{invalid`))
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/register", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Erro - Falha ao registrar", func(t *testing.T) {
		mockUseCase := new(mocks.MockRegisterUseCase)
		mockUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("erro no banco"))

		handler := handlers.NewRegisterHandler(mockUseCase, nil)

		reqBody := `{"name": "John", "email": "john@example.com", "password": "123"}`
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/register", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Erro - Formato de e-mail inválido", func(t *testing.T) {
		mockUseCase := new(mocks.MockRegisterUseCase)
		handler := handlers.NewRegisterHandler(mockUseCase, nil)

		mockUseCase.AssertNotCalled(t, "Execute")

		reqBody := `{
			"name": "John Doe",
			"email": "email-invalido",
			"password": "senha123"
		}`
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/register", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"invalid email format"}`, resp.Body.String())
	})

	t.Run("Erro - Falha ao buscar usuário após registro", func(t *testing.T) {
		mockUseCase := new(mocks.MockRegisterUseCase)
		mockRepo := new(mocks.MockUserRepo)

		handler := handlers.NewRegisterHandler(mockUseCase, mockRepo)

		mockUseCase.On("Execute", mock.Anything, dto.RegisterParams{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "senha123",
		}).Return(nil)

		email, _ := vo.NewEmail("john@example.com")
		mockRepo.On("GetByEmail", mock.Anything, email).
			Return(nil, errors.New("erro no banco"))

		reqBody := `{
			"name": "John Doe",
			"email": "john@example.com",
			"password": "senha123"
	}`
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/register", handler.Handle)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"failed to fetch user"}`, resp.Body.String())
		mockUseCase.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
