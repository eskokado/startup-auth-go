package handlers

import (
	"net/http"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/gin-gonic/gin"
)

type ForgotPasswordHandler struct {
	useCase usecase.RequestPasswordResetInterface
}

func NewForgotPasswordHandler(uc usecase.RequestPasswordResetInterface) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{useCase: uc}
}

func (h *ForgotPasswordHandler) Handle(c *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := vo.NewEmail(input.Email)
	if err != nil {
		// Retorna o erro específico em vez de sempre "invalid email format"
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.Execute(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
