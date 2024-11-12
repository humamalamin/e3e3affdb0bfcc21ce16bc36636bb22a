package repository

import (
	"context"
	"errors"
	"fmt"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategoryByID(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

// DeleteCategoryByID implements CategoryRepository.
func (c *categoryRepository) DeleteCategoryByID(ctx context.Context, id int64) error {
	var modelCategory model.Category

	err = c.db.Where("id =?", id).Preload("Contents").First(&modelCategory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		code = "[Repository] DeleteCategoryByID - 1"
		log.Errorw(code, err)
		return err
	}

	if len(modelCategory.Contents) > 0 {
		code = "[Repository] DeleteCategoryByID - 2"
		err = errors.New("category cannot be deleted because it has related contents")
		log.Errorw(code, err)
		return err
	}
	err = c.db.Where("id =?", id).Delete(&modelCategory).Error
	if err != nil {
		code = "[Repository] DeleteCategoryByID - 3"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// EditCategoryByID implements CategoryRepository.
func (c *categoryRepository) EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err := c.db.Table("categories").Where("slug =?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[Repository] EditCategoryByID - 1"
		log.Errorw(code, err)
		return err
	}

	countSlug = countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)

	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(modelCategory).Error
	if err != nil {
		code = "[Repository] EditCategoryByID - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err := c.db.Table("categories").Where("slug =?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[Repository] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	countSlug = countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)

	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Create(&modelCategory).Error
	if err != nil {
		code = "[Repository] CreateCategory - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetCategoryByID implements CategoryRepository.
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	var modelCategory model.Category
	err := c.db.Where("id =?", id).Preload("User").First(&modelCategory).Error
	if err != nil {
		code = "[Repository] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return &entity.CategoryEntity{
		ID:    modelCategory.ID,
		Title: modelCategory.Title,
		Slug:  modelCategory.Slug,
		User: entity.UserEntity{
			ID:    modelCategory.CreatedByID,
			Name:  modelCategory.User.Name,
			Email: modelCategory.User.Email,
		},
	}, nil
}

// GetCategories implements CategoryRepository.
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCategories []model.Category

	err = c.db.Order("created_at DESC").Preload("User").Find(&modelCategories).Error
	if err != nil {
		code = "[Repository] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCategories) == 0 {
		code = "[Repository] GetCategories - 2"
		err = errors.New("data not found")
		log.Errorw(code, err)
		return nil, err
	}

	var resp []entity.CategoryEntity
	for _, val := range modelCategories {
		resp = append(resp, entity.CategoryEntity{
			ID:    val.ID,
			Title: val.Title,
			Slug:  val.Slug,
			User: entity.UserEntity{
				ID:    val.CreatedByID,
				Name:  val.User.Name,
				Email: val.User.Email,
			},
		})
	}

	return resp, nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
