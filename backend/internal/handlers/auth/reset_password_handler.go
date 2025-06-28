package handlers

import (
	"errors"
	"net/http"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/gin-gonic/gin"
)

type ResetPasswordHandler struct {
	useCase usecase.ResetPasswordInterface
}

func NewResetPasswordHandler(uc usecase.ResetPasswordInterface) *ResetPasswordHandler {
	return &ResetPasswordHandler{useCase: uc}
}

func (h *ResetPasswordHandler) Handle(c *gin.Context) {
	var input dto.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.useCase.Execute(c.Request.Context(), input.Token, input.Password); err != nil {
		switch {
		case errors.Is(err, msgerror.AnErrInvalidToken),
			errors.Is(err, msgerror.AnErrExpiredToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		}
		return
	}

	// Garante que nenhum conteúdo seja retornado
	c.Status(http.StatusNoContent)
}
