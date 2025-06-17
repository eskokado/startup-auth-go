package usecase_test

import (
	"context"
	"errors"
	"testing"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var originalWithName = (*entity.User).WithName

func TestUpdateNameUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	validUserID := vo.NewID()
	validUser := &entity.User{ID: validUserID}
	oldName := "Existing Name"
	oldNameVo, _ := vo.NewName(oldName, 3, 50)
	userWithName := &entity.User{
		ID:   validUserID,
		Name: oldNameVo,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockRepo.On("Save", ctx, mock.Anything).Return(validUser, nil)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "Novo Nome Diferente")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return((*entity.User)(nil), msgerror.AnErrNotFound)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "New Name")

		assert.ErrorIs(t, err, msgerror.AnErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UnexpectedGetError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return((*entity.User)(nil), errors.New("db connection failed"))

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "New Name")

		assert.ErrorContains(t, err, "failed to get user")
		assert.ErrorContains(t, err, "db connection failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("NilUserAfterGet", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return((*entity.User)(nil), nil)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "New Name")

		assert.ErrorIs(t, err, msgerror.AnErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidName", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "ab") // Nome muito curto

		assert.ErrorIs(t, err, msgerror.AnErrInvalidName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SaveError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockRepo.On("Save", ctx, mock.Anything).Return((*entity.User)(nil), errors.New("save failed"))

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "Valid Name")

		assert.ErrorContains(t, err, "failed to save user")
		mockRepo.AssertExpectations(t)
	})

	t.Run("SameName", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(userWithName, nil)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, oldName) // Mesmo nome atual

		assert.ErrorIs(t, err, msgerror.AnErrNameDifferent)
		mockRepo.AssertNotCalled(t, "Save")
	})
}
