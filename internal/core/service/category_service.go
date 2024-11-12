package service

import (
	"context"
	"latihan-portal-news/internal/adapter/repository"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategoryByID(ctx context.Context, id int64) error
}

type categoryService struct {
	repository repository.CategoryRepository
}

// DeleteCategoryByID implements CategoryService.
func (c *categoryService) DeleteCategoryByID(ctx context.Context, id int64) error {
	err = c.repository.DeleteCategoryByID(ctx, id)
	if err != nil {
		code = "[Service] DeleteCategoryByID - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// EditCategoryByID implements CategoryService.
func (c *categoryService) EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug
	err := c.repository.EditCategoryByID(ctx, req)
	if err != nil {
		code = "[Service] EditCategoryByID - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug
	err = c.repository.CreateCategory(ctx, req)
	if err != nil {
		code = "[Service] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetCategoryByID implements CategoryService.
func (c *categoryService) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	result, err := c.repository.GetCategoryByID(ctx, id)
	if err != nil {
		code = "[Service] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

// GetCategories implements CategoryService.
func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	results, err := c.repository.GetCategories(ctx)
	if err != nil {
		code = "[Service] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repository: repository}
}
