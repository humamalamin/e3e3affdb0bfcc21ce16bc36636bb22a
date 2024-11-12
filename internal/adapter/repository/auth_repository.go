package repository

import (
	"context"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var err error
var code string

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, req entity.RequestLogin) (*entity.UserEntity, error)
}

type authRepository struct {
	db *gorm.DB
}

// GetUserByEmail implements AuthRepository.
func (a *authRepository) GetUserByEmail(ctx context.Context, req entity.RequestLogin) (*entity.UserEntity, error) {
	var modelUser model.User

	err = a.db.Where("email = ?", req.Email).First(&modelUser).Error
	if err != nil {
		code = "[Repository] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.UserEntity{
		ID:       modelUser.ID,
		Name:     modelUser.Name,
		Email:    modelUser.Email,
		Password: modelUser.Password,
	}

	return &resp, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
