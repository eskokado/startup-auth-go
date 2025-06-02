package handlers

import (
	"net/http"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/port/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginHandler struct {
	loginUseCase  usecase.LoginInterface
	tokenProvider providers.TokenProvider
}

func NewLoginHandler(loginUseCase usecase.LoginInterface, tokenProvider providers.TokenProvider) *LoginHandler {
	return &LoginHandler{
		loginUseCase:  loginUseCase,
		tokenProvider: tokenProvider,
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

	// Gerar token JWT
	claims := providers.Claims{
		UserID: loginResult.UserID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   loginResult.Email.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Expira em 24h,
		},
	}
	token, err := h.tokenProvider.Generate(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	output := dto.LoginOutput{AccessToken: token}
	c.JSON(http.StatusOK, output)
}
