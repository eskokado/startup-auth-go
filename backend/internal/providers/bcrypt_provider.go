package providers

import (
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"golang.org/x/crypto/bcrypt"
)

type BcryptProvider struct {
	cost int
}

func NewBcryptProvider(cost int) *BcryptProvider {
	return &BcryptProvider{cost: cost}
}

func (b *BcryptProvider) Encrypt(password string) (string, error) {
	if password == "" {
		return "", msgerror.AnErrEmptyPassword
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(hash), nil
}

func (b *BcryptProvider) Compare(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
