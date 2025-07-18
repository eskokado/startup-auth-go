package middleware

import (
	"strings"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(
	tokenProvider providers.TokenProvider,
	blacklistProvider providers.BlacklistProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extrair o token do header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		// 2. Verificar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := parts[1]

		// 3. Verificar e o token está na blacklist
		keyToken := "startup-auth-go:" + tokenString + ":Token"
		blacklisted, err := blacklistProvider.ExistsKey(c.Request.Context(), keyToken)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			return
		}
		if !blacklisted {
			c.AbortWithStatusJSON(401, gin.H{"error": "token revoked"})
			return
		}

		// 4. Validar token e obter claims
		rawClaims, err := tokenProvider.Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
			return
		}

		// 5. Converter claims para o tipo correto
		claims, ok := rawClaims.(providers.Claims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims structure"})
			return
		}

		// 6. Extrair userID das claims
		userID := claims.UserID
		if userID == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "UserID not found in token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
