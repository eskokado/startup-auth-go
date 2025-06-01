package entity

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

var GenerateSecureToken = generateSecureToken

type User struct {
	ID                   vo.ID
	Name                 vo.Name
	Email                vo.Email
	PasswordHash         vo.PasswordHash
	ImageURL             vo.URL
	CreatedAt            time.Time
	PasswordResetToken   string
	PasswordResetExpires time.Time
}

func NewUser(
	id vo.ID,
	name vo.Name,
	email vo.Email,
	passwordHash vo.PasswordHash,
	imageURL *vo.URL,
) (*User, error) {

	if email.IsEmpty() {
		return nil, msgerror.AnErrEmptyEmail
	}

	if passwordHash.IsEmpty() {
		return nil, msgerror.AnErrWeakPassword
	}

	user := &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if imageURL != nil {
		user.ImageURL = *imageURL
	}

	return user, nil
}

func CreateUser(
	name string,
	email string,
	password string,
	imageURL string,
) (*User, error) {
	id := vo.NewID()

	validName, err := vo.NewName(name, 3, 50)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", msgerror.AnErrInvalidUser, err)
	}

	validEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", msgerror.AnErrInvalidUser, err)
	}

	passwordHash, err := vo.NewPasswordHash(password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", msgerror.AnErrInvalidUser, err)
	}

	var url *vo.URL
	if imageURL != "" {
		u, err := vo.NewURL(imageURL)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", msgerror.AnErrInvalidUser, err)
		}
		url = &u
	}

	return NewUser(id, validName, validEmail, passwordHash, url)
}

func (u *User) WithName(newName vo.Name) (*User, error) {
	if newName.IsEmpty() {
		return nil, msgerror.AnErrInvalidName
	}

	return &User{
		ID:           u.ID,
		Name:         newName,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		ImageURL:     u.ImageURL,
		CreatedAt:    u.CreatedAt,
	}, nil
}

func (u *User) WithPasswordHash(newHash vo.PasswordHash) (*User, error) {
	if newHash.IsEmpty() {
		return nil, msgerror.AnErrWeakPassword
	}
	return &User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: newHash,
		ImageURL:     u.ImageURL,
		CreatedAt:    u.CreatedAt,
	}, nil
}

func (u *User) Equal(other *User) bool {
	return u.ID.Equal(other.ID)
}

func (u *User) VerifyPassword(password string) bool {
	return u.PasswordHash.Verify(password)
}
func (u *User) GeneratePasswordResetToken() error {
	token, err := GenerateSecureToken()
	if err != nil {
		return err
	}
	u.PasswordResetToken = token
	u.PasswordResetExpires = time.Now().Add(1 * time.Hour)
	return nil
}

func (u *User) ClearResetToken() {
	u.PasswordResetToken = ""
	u.PasswordResetExpires = time.Time{}
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
