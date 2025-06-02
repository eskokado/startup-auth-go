package handlers

import (
	"net/http"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	registerUseCase usecase.RegisterInterface
	userRepo        repository.UserRepository
}

func NewRegisterHandler(
	registerUseCase usecase.RegisterInterface,
	userRepo repository.UserRepository,
) *RegisterHandler {
	return &RegisterHandler{
		registerUseCase: registerUseCase,
		userRepo:        userRepo,
	}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var input dto.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validação do email ANTES do caso de uso
	email, err := vo.NewEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	params := dto.RegisterParams{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.registerUseCase.Execute(c.Request.Context(), params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	user, err := h.userRepo.GetByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	output := dto.RegisterUserOutput{
		ID:       user.ID.String(),
		Name:     user.Name.String(),
		Email:    user.Email.String(),
		ImageURL: user.ImageURL.String(),
	}
	c.JSON(http.StatusCreated, output)
}
