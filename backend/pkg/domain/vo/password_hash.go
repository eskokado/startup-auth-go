package vo

import (
	"regexp"

	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash struct {
	value string
}

func NewPasswordHash(password string) (PasswordHash, error) {
	if !isBcryptHash(password) {
		if len(password) < 8 {
			return PasswordHash{}, msgerror.AnErrPasswordInvalid
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return PasswordHash{}, err
		}
		return PasswordHash{value: string(hash)}, nil
	}

	return PasswordHash{value: password}, nil
}

func (ph PasswordHash) String() string {
	return ph.value
}

func (ph PasswordHash) IsEmpty() bool {
	return ph.value == ""
}

func (ph PasswordHash) Verify(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ph.value), []byte(password))
	return err == nil
}

// Função melhorada com regex para validação precisa
func isBcryptHash(s string) bool {
	// Regex para validar formato bcrypt
	bcryptRegex := regexp.MustCompile(`^\$2[ayb]\$[0-9]{2}\$[./A-Za-z0-9]{53}$`)
	return bcryptRegex.MatchString(s)
}

func IsBcryptHash(s string) bool {
	return isBcryptHash(s)
}
