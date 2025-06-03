package handlers

import (
	"net/http"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/gin-gonic/gin"
)

type UpdatePasswordHandler struct {
	updatePasswordUseCase usecase.UpdatePasswordInterface
}

func NewUpdatePasswordHandler(updatePasswordUseCase usecase.UpdatePasswordInterface) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{
		updatePasswordUseCase: updatePasswordUseCase,
	}
}

func (h *UpdatePasswordHandler) Handle(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := vo.ParseID(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": msgerror.AnErrInvalidID.Error()})
		return
	}

	var input dto.UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.updatePasswordUseCase.Execute(c.Request.Context(), userID, input.CurrentPassword, input.NewPassword); err != nil {
		switch err {
		case msgerror.AnErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case msgerror.AnErrWeakPassword, msgerror.AnErrPasswordInvalid:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case msgerror.AnErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		}
		return
	}

	c.Status(http.StatusOK)
}
