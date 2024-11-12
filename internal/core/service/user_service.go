package service

import (
	"context"
	"latihan-portal-news/internal/adapter/repository"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	UpdatePassword(ctx context.Context, newPassword string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	repo repository.UserRepository
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		code = "[Service] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return user, nil
}

// UpdatePassword implements UserService.
func (u *userService) UpdatePassword(ctx context.Context, newPassword string, id int64) error {
	password, err := conv.HashPassword(newPassword)
	if err != nil {
		code = "[Service] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	err = u.repo.UpdatePassword(ctx, password, id)
	if err != nil {
		code = "[Service] UpdatePassword - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
