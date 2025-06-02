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

func TestUpdateNameUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	validUserID := vo.NewID()
	validUser := &entity.User{ID: validUserID}

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
		mockRepo.On("GetByID", ctx, validUserID).Return(&entity.User{}, msgerror.AnErrUserNotFound)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "New Name")

		assert.ErrorContains(t, err, "failed to get user")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidName", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "ab") // Nome muito curto

		assert.ErrorIs(t, err, msgerror.AnErrInvalidUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SaveError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockRepo.On("Save", ctx, mock.Anything).Return(&entity.User{}, errors.New("save failed"))

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err := uc.Execute(ctx, validUserID, "Valid Name")

		assert.ErrorContains(t, err, "save failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("SameName", func(t *testing.T) {
		oldName := "Existing Name"
		oldNameVo, err := vo.NewName(oldName, 3, 50)
		assert.NoError(t, err)
		userWithName := &entity.User{
			ID:   validUserID,
			Name: oldNameVo,
		}

		mockRepo := new(mocks.MockUserRepo)
		mockRepo.On("GetByID", ctx, validUserID).Return(userWithName, nil)

		uc := usecase.NewUpdateNameUseCase(mockRepo)
		err = uc.Execute(ctx, validUserID, oldName) // Mesmo nome atual

		assert.Error(t, err)
		assert.ErrorContains(t, err, "new name must be different") // Mensagem específica

		mockRepo.AssertNotCalled(t, "Save")
	})
}
