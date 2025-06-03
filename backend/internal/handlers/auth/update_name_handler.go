package handlers

import (
	"net/http"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/gin-gonic/gin"
)

type UpdateNameHandler struct {
	updateNameUseCase usecase.UpdateNameInterface
}

func NewUpdateNameHandler(updateNameUseCase usecase.UpdateNameInterface) *UpdateNameHandler {
	return &UpdateNameHandler{
		updateNameUseCase: updateNameUseCase,
	}
}

func (h *UpdateNameHandler) Handle(c *gin.Context) {
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

	var input dto.UpdateNameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.updateNameUseCase.Execute(c.Request.Context(), userID, input.Name); err != nil {
		switch err {
		case msgerror.AnErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case msgerror.AnErrInvalidName, msgerror.AnErrNameTooShort, msgerror.AnErrNameTooLong:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update name"})
		}
		return
	}

	c.Status(http.StatusOK)
}
