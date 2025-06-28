package handlers

import (
	"net/http"
	"strings"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	logoutUseCase usecase.LogoutInterface
}

func NewLogoutHandler(logoutUseCase usecase.LogoutInterface) *LogoutHandler {
	return &LogoutHandler{
		logoutUseCase: logoutUseCase,
	}
}

func (h *LogoutHandler) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
		return
	}

	token := parts[1]

	err := h.logoutUseCase.Execute(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.Status(http.StatusOK)
}
