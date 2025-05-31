package repository

import (
	"context"
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"gorm.io/gorm"
)

type GormUser struct {
	ID                   string    `gorm:"primaryKey;type:varchar(36)"`
	Name                 string    `gorm:"type:varchar(100);not null"`
	Email                string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash         string    `gorm:"type:varchar(255);not null"`
	ImageURL             string    `gorm:"type:varchar(255)"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	PasswordResetToken   string    `gorm:"type:varchar(255)"`
	PasswordResetExpires time.Time `gorm:"type:datetime"`
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) toDBModel(user *entity.User) *GormUser {
	return &GormUser{
		ID:                   user.ID.String(),
		Name:                 user.Name.String(),
		Email:                user.Email.String(),
		PasswordHash:         user.PasswordHash.String(),
		ImageURL:             user.ImageURL.String(),
		PasswordResetToken:   user.PasswordResetToken,
		PasswordResetExpires: user.PasswordResetExpires,
	}
}

func (r *GormUserRepository) fromDBModel(dbUser *GormUser) (*entity.User, error) {
	id, err := vo.ParseID(dbUser.ID)
	if err != nil {
		return nil, err
	}

	name, err := vo.NewName(dbUser.Name, 3, 50)
	if err != nil {
		return nil, err
	}

	email, err := vo.NewEmail(dbUser.Email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := vo.NewPasswordHash(dbUser.PasswordHash)
	if err != nil {
		return nil, err
	}

	imageURL, err := vo.NewURL(dbUser.ImageURL)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                   id,
		Name:                 name,
		Email:                email,
		PasswordHash:         passwordHash,
		ImageURL:             imageURL,
		PasswordResetToken:   dbUser.PasswordResetToken,
		PasswordResetExpires: dbUser.PasswordResetExpires,
	}, nil
}

func (r *GormUserRepository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser := r.toDBModel(user)

	result := r.db.WithContext(ctx).Save(dbUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.fromDBModel(dbUser)
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	var dbUser GormUser
	result := r.db.WithContext(ctx).Where("email = ?", email.String()).First(&dbUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return r.fromDBModel(&dbUser)
}

func (r *GormUserRepository) GetByResetToken(ctx context.Context, token string) (*entity.User, error) {
	var dbUser GormUser
	result := r.db.WithContext(ctx).Where("password_reset_token = ?", token).First(&dbUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return r.fromDBModel(&dbUser)
}

func (r *GormUserRepository) GetByID(ctx context.Context, id vo.ID) (*entity.User, error) {
	var dbUser GormUser
	result := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&dbUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return r.fromDBModel(&dbUser)
}
