package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/golang-jwt/jwt/v5"
)

type LoginUsecase struct {
	userRepo          repository.UserRepository
	cryptoProvider    providers.CryptoProvider
	tokenProvider     providers.TokenProvider
	blacklistProvider providers.BlacklistProvider
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	cryptoProvider providers.CryptoProvider,
	tokenProvider providers.TokenProvider,
	blacklistProvider providers.BlacklistProvider,
) *LoginUsecase {
	return &LoginUsecase{
		userRepo:          userRepo,
		cryptoProvider:    cryptoProvider,
		tokenProvider:     tokenProvider,
		blacklistProvider: blacklistProvider,
	}
}
func (h *LoginUsecase) Execute(ctx context.Context, email string, password string) (dto.LoginResult, error) {
	validationErrs := msgerror.NewValidationErrors()

	// Validação básica de campos
	if email == "" {
		validationErrs.Add("email", "cannot be empty")
	} else if _, err := vo.NewEmail(email); err != nil {
		validationErrs.Add("email", err.Error())
	}

	if password == "" {
		validationErrs.Add("password", "cannot be empty")
	} else if len(password) < 8 {
		validationErrs.Add("password", "must be at least 8 characters")
	}

	// Se houver erros de validação básicos, retorne imediatamente
	if validationErrs.HasErrors() {
		return dto.LoginResult{}, validationErrs
	}

	// Convertemos para vo.Email (já validado acima, então não haverá erro)
	validEmail, _ := vo.NewEmail(email)

	user, err := h.userRepo.GetByEmail(ctx, validEmail)
	if errors.Is(err, msgerror.AnErrNotFound) {
		// Por segurança, não revelamos que o usuário não existe
		return dto.LoginResult{}, msgerror.AnErrInvalidCredentials
	}
	if err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to get user", err)
	}

	match, err := h.cryptoProvider.Compare(password, user.PasswordHash.String())
	if err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to verify password", err)
	}
	if !match {
		return dto.LoginResult{}, msgerror.AnErrInvalidCredentials
	}

	// Gerar token JWT
	claims := providers.Claims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token, err := h.tokenProvider.Generate(claims)
	if err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to generate token", err)
	}

	// Salvar dados no Redis com prefixo
	prefix := "startup-auth-go"
	ttl := 24 * time.Hour

	// UserID
	keyUserID := prefix + ":UserID"
	if err := h.blacklistProvider.SetWithKey(ctx, keyUserID, user.ID.String(), ttl); err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to save UserID", err)
	}

	// Name
	keyName := prefix + ":Name"
	if err := h.blacklistProvider.SetWithKey(ctx, keyName, user.Name.String(), ttl); err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to save Name", err)
	}

	// Email
	keyEmail := prefix + ":Email"
	if err := h.blacklistProvider.SetWithKey(ctx, keyEmail, user.Email.String(), ttl); err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to save Email", err)
	}

	// Token
	keyToken := prefix + ":Token"
	if err := h.blacklistProvider.SetWithKey(ctx, keyToken, token, ttl); err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to save Token", err)
	}

	// CreatedAt (convertido para string)
	keyCreatedAt := prefix + ":CreatedAt"
	createdAtStr := user.CreatedAt.Format(time.RFC3339)
	if err := h.blacklistProvider.SetWithKey(ctx, keyCreatedAt, createdAtStr, ttl); err != nil {
		return dto.LoginResult{}, msgerror.Wrap("failed to save CreatedAt", err)
	}

	return dto.LoginResult{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Token:     token,
	}, nil
}
