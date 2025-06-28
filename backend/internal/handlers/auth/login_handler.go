package handlers

import (
	"net/http"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUseCase usecase.LoginInterface
}

func NewLoginHandler(loginUseCase usecase.LoginInterface) *LoginHandler {
	return &LoginHandler{
		loginUseCase: loginUseCase,
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	loginResult, err := h.loginUseCase.Execute(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	output := dto.LoginOutput{
		AccessToken: loginResult.Token,
		User: dto.UserOutput{
			Id:    loginResult.UserID.String(),
			Name:  loginResult.Name.String(),
			Email: loginResult.Email.String(),
		},
	}
	c.JSON(http.StatusOK, output)
}
