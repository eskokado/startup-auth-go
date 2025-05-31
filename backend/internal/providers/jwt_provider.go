package provider

import (
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTProvider(secret string, expiry time.Duration) *JWTProvider {
	return &JWTProvider{
		secretKey: []byte(secret),
		expiry:    expiry,
	}
}

func (j *JWTProvider) Generate(claims interface{}) (string, error) {
	c, ok := claims.(providers.Claims)
	if !ok {
		return "", errors.New("tipo de claims inválido")
	}

	expirationTime := time.Now().Add(j.expiry)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     c.RegisteredClaims.Subject,
		"exp":     expirationTime.Unix(),
		"user_id": c.UserID,
	})
	return token.SignedString(j.secretKey)
}

func (j *JWTProvider) Validate(token string) (interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	subject, ok := claims["sub"].(string)
	if !ok {
		return providers.Claims{}, errors.New("subject não encontrado ou inválido")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return providers.Claims{}, errors.New("user_id não encontrado ou inválido")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return providers.Claims{}, errors.New("exp não encontrado ou inválido")
	}
	expirationTime := time.Unix(int64(expFloat), 0)

	return providers.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}, nil
}
