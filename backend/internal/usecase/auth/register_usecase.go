package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type RegisterUsecase struct {
	userRepo       repository.UserRepository
	cryptoProvider providers.CryptoProvider
}

func NewRegisterUsecase(
	userRepo repository.UserRepository,
	cryptoProvider providers.CryptoProvider,
) *RegisterUsecase {
	return &RegisterUsecase{
		userRepo:       userRepo,
		cryptoProvider: cryptoProvider,
	}
}

func (h *RegisterUsecase) Execute(ctx context.Context, input dto.RegisterParams) error {
	validationErrs := msgerror.NewValidationErrors()

	// Validação básica de campos
	if input.Password != input.PasswordConfirmation {
		validationErrs.Add("password_confirmation", "passwords do not match")
	}
	if len(input.Password) < 8 {
		validationErrs.Add("password", "must be at least 8 characters")
	}

	// Validação de objetos de valor
	name, nameErr := vo.NewName(input.Name, 3, 100)
	if nameErr != nil {
		validationErrs.Add("name", nameErr.Error())
	}

	email, emailErr := vo.NewEmail(input.Email)
	if emailErr != nil {
		validationErrs.Add("email", emailErr.Error())
	}

	var imageURL vo.URL
	if input.ImageURL != "" {
		url, urlErr := vo.NewURL(input.ImageURL)
		if urlErr != nil {
			validationErrs.Add("image_url", urlErr.Error())
		} else {
			imageURL = url
		}
	}

	// Se houver erros de validação, retornar imediatamente
	if validationErrs.HasErrors() {
		return validationErrs
	}

	// Verificação de e-mail único (só se o e-mail for válido)
	if emailErr == nil {
		existingUser, err := h.userRepo.GetByEmail(ctx, email)
		if err != nil && !errors.Is(err, msgerror.AnErrNotFound) {
			// Erro de infraestrutura, não de validação
			return msgerror.Wrap("failed to check email existence", err)
		}
		if existingUser != nil {
			validationErrs.Add("email", msgerror.AnErrUserExists.Error())
			return validationErrs
		}
	}

	// Criptografia de senha
	hashedPassword, err := h.cryptoProvider.Encrypt(input.Password)
	if err != nil {
		return msgerror.Wrap("failed to secure password", err)
	}

	passwordHashed, err := vo.NewPasswordHash(hashedPassword)
	if err != nil {
		return msgerror.Wrap("failed to create password hash", err)
	}

	// Criação do usuário
	newUser := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHashed,
		ImageURL:     imageURL,
		CreatedAt:    time.Now(),
	}

	// Persistência
	savedUser, err := h.userRepo.Save(ctx, newUser)
	if err != nil {
		return msgerror.Wrap("failed to create user", err)
	}

	if savedUser == nil {
		return msgerror.AnErrNoSavedUser
	}

	return nil
}
