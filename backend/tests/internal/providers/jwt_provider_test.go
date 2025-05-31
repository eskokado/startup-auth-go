package providers_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	auth "github.com/eskokado/startup-auth-go/backend/internal/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTProvider(t *testing.T) {
	provider := auth.NewJWTProvider("test-secret-key", 15*time.Minute)
	userID := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("Generate and Validate Token", func(t *testing.T) {
		claims := providers.Claims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID,
			},
		}
		token, err := provider.Generate(claims)
		if err != nil {
			t.Fatalf("Token generation failed: %v", err)
		}

		validatedClaims, err := provider.Validate(token)
		if err != nil {
			t.Fatalf("Token validation failed: %v", err)
		}

		vc, ok := validatedClaims.(providers.Claims)
		if !ok {
			t.Fatal("Validated claims are not of type providers.Claims")
		}
		if vc.UserID != userID {
			t.Errorf("Expected UserID %s, got %s", userID, vc.UserID)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		expiredProvider := auth.NewJWTProvider("test-secret-key", -5*time.Minute)
		claims := providers.Claims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID,
			},
		}
		token, _ := expiredProvider.Generate(claims)

		_, err := provider.Validate(token)
		if err == nil {
			t.Error("Expected expired token error")
		} else {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return
			}
			t.Errorf("Expected token expiration error, got: %v", err)
		}
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		claims := providers.Claims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID,
			},
		}
		token, _ := provider.Generate(claims)
		invalidProvider := auth.NewJWTProvider("wrong-secret-key", 15*time.Minute)

		_, err := invalidProvider.Validate(token)
		if err == nil {
			t.Error("Expected signature validation error")
		}
	})

	t.Run("Generate with Invalid Claims Type", func(t *testing.T) {
		invalidClaims := "not_a_struct"
		_, err := provider.Generate(invalidClaims)
		if err == nil {
			t.Error("Expected error for invalid claims type")
		} else if err.Error() != "tipo de claims inválido" {
			t.Errorf("Expected 'tipo de claims inválido', got: %v", err)
		}
	})

	t.Run("Malformed Token", func(t *testing.T) {
		_, err := provider.Validate("invalid.token.string")
		if err == nil {
			t.Error("Expected error for malformed token")
		}
	})

	t.Run("Missing Subject Claim", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
			"user_id": userID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		_, err := provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for missing subject")
		} else if err.Error() != "subject não encontrado ou inválido" {
			t.Errorf("Expected 'subject não encontrado ou inválido', got: %v", err)
		}
	})

	t.Run("Invalid Subject Claim Type", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":     123,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
			"user_id": userID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		_, err := provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for invalid subject type")
		} else if err.Error() == "subject não encontrado ou inválido" {
			// Aceita o erro personalizado
		} else {
			t.Errorf("Expected 'subject não encontrado ou inválido', got: %v", err)
		}
	})

	t.Run("Missing UserID Claim", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(15 * time.Minute).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		_, err := provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for missing user_id")
		} else if err.Error() != "user_id não encontrado ou inválido" {
			t.Errorf("Expected 'user_id não encontrado ou inválido', got: %v", err)
		}
	})

	t.Run("Invalid UserID Claim Type", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":     userID,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
			"user_id": 12345,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		_, err := provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for invalid user_id type")
		} else if err.Error() == "user_id não encontrado ou inválido" {
			// Aceita o erro personalizado
		} else {
			t.Errorf("Expected 'user_id não encontrado ou inválido', got: %v", err)
		}
	})

	t.Run("Missing Expiration Claim", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":     userID,
			"user_id": userID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		_, err := provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for missing exp")
		} else if err.Error() != "exp não encontrado ou inválido" {
			t.Errorf("Expected 'exp não encontrado ou inválido', got: %v", err)
		}
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		_, err := provider.Validate("invalid-token-format")
		if err == nil {
			t.Error("Expected error for invalid token format")
		}
	})

	t.Run("Invalid Token Signature Format", func(t *testing.T) {
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		parts := strings.Split(validToken, ".")
		invalidToken := parts[0] + "." + parts[1] + ".invalid_signature"

		_, err := provider.Validate(invalidToken)
		if err == nil {
			t.Error("Expected signature validation error")
		}
	})

	t.Run("Unsupported Signing Method", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":     userID,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
			"user_id": userID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret-key"))

		_, err := provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for unsupported signing method")
		}
	})

	t.Run("Invalid Token Structure", func(t *testing.T) {
		_, err := provider.Validate("header.claims")
		if err == nil {
			t.Error("Expected error for invalid token structure")
		}
	})

	t.Run("Empty Token", func(t *testing.T) {
		_, err := provider.Validate("")
		if err == nil {
			t.Error("Expected error for empty token")
		}
	})

	t.Run("Nil Claims", func(t *testing.T) {
		_, err := provider.Generate(nil)
		if err == nil {
			t.Error("Expected error for nil claims")
		}
	})

	t.Run("Token Not Valid", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":     userID,
			"exp":     "invalid_expiration", // valor inválido para exp
			"user_id": userID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test-secret-key"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for invalid token")
		} else if !strings.Contains(err.Error(), "token") {
			t.Errorf("Expected token error, got: %v", err)
		}
	})

	t.Run("Invalid Claims Type", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)

		token.Claims = jwt.RegisteredClaims{}

		tokenString, err := token.SignedString([]byte("test-secret-key"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for invalid claims type")
		} else if err.Error() != "subject não encontrado ou inválido" {
			t.Errorf("Expected 'subject não encontrado ou inválido', got: %v", err)
		}
	})
}
