package providers_test

import (
	"testing"

	crypto "github.com/eskokado/startup-auth-go/backend/internal/providers"
)

func TestBcryptProvider(t *testing.T) {
	provider := crypto.NewBcryptProvider(12)

	t.Run("Encrypt and Compare", func(t *testing.T) {
		password := "SecurePass123!"

		hash, err := provider.Encrypt(password)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		match, err := provider.Compare(password, hash)
		if err != nil || !match {
			t.Errorf("Password comparison failed: %v", err)
		}
	})

	t.Run("Compare Invalid Password", func(t *testing.T) {
		hash, _ := provider.Encrypt("GoodPass123!")

		match, err := provider.Compare("WrongPass456@", hash)
		if err == nil || match {
			t.Error("Expected password mismatch")
		}
	})

	t.Run("Empty Password", func(t *testing.T) {
		_, err := provider.Encrypt("")
		if err == nil {
			t.Error("Expected error for empty password")
		}
	})
}
