package repository

import (
	"context"
	"fmt"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/domain/model"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ContentRepository interface {
	GetContents(ctx context.Context, query *entity.QueryString) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContentByID(ctx context.Context, id int64) error

	GetContentByCategoryID(ctx context.Context, categoryID int64, search string) ([]entity.ContentEntity, error)
}

type contentRepository struct {
	db *gorm.DB
}

// GetContentByCategoryID implements ContentRepository.
func (c *contentRepository) GetContentByCategoryID(ctx context.Context, categoryID int64, search string) ([]entity.ContentEntity, error) {
	var modelContents []model.Content

	if search == "" {
		err = c.db.Where("category_id =?", categoryID).Preload("User", "Category").Find(&modelContents).Error
	} else {
		err = c.db.Where("category_id =? AND (title ilike ? OR description ilike ? OR excerpt ilike ? OR tags ilike ?",
			categoryID, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").
			Preload("User", "Category").Find(&modelContents).Error
	}
	if err != nil {
		code = "[Repository] GetContentByCategoryID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			CategoryID:  val.CategoryID,
			CreatedByID: val.CreatedByID,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Image:       val.Image,
			Status:      val.Status,
			Tags:        tags,
			Category: entity.CategoryEntity{
				ID:    val.Category.ID,
				Title: val.Category.Title,
				Slug:  val.Category.Slug,
			},
			User: entity.UserEntity{
				ID:   val.User.ID,
				Name: val.User.Name,
			},
			CreatedAt: val.CreatedAt,
		}

		resps = append(resps, resp)
	}

	return resps, nil
}

// DeleteContentByID implements ContentRepository.
func (c *contentRepository) DeleteContentByID(ctx context.Context, id int64) error {
	err = c.db.Where("id =?", id).Delete(&model.Content{}).Error
	if err != nil {
		code = "[Repository] DeleteContentByID - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Status:      req.Status,
		Tags:        tags,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelContent).Error
	if err != nil {
		code = "[Repository] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// CreateContent implements ContentRepository.
func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Status:      req.Status,
		Tags:        tags,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code = "[Repository] CreateContent - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentRepository.
func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content

	err = c.db.Where("id = ?", id).Preload("User", "Category").First(&modelContent).Error
	if err != nil {
		code = "[Repository] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	tags := strings.Split(modelContent.Tags, ",")
	resp := entity.ContentEntity{
		ID:          id,
		Title:       modelContent.Title,
		CategoryID:  modelContent.CategoryID,
		CreatedByID: modelContent.CreatedByID,
		Excerpt:     modelContent.Excerpt,
		Description: modelContent.Description,
		Image:       modelContent.Image,
		Status:      modelContent.Status,
		Tags:        tags,
		Category: entity.CategoryEntity{
			ID:    modelContent.Category.ID,
			Title: modelContent.Category.Title,
			Slug:  modelContent.Category.Slug,
		},
		User: entity.UserEntity{
			ID:   modelContent.User.ID,
			Name: modelContent.User.Name,
		},
		CreatedAt: modelContent.CreatedAt,
	}

	return &resp, nil
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context, query *entity.QueryString) ([]entity.ContentEntity, error) {
	var modelContents []model.Content

	order := fmt.Sprintf("%s %s", query.Order, query.OrderType)

	offset := (query.Page - 1) * query.Limit
	err = c.db.Preload("User", "Category").
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelContents).Error
	if err != nil {
		code = "[Repository] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			CategoryID:  val.CategoryID,
			CreatedByID: val.CreatedByID,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Image:       val.Image,
			Status:      val.Status,
			Tags:        tags,
			Category: entity.CategoryEntity{
				ID:    val.Category.ID,
				Title: val.Category.Title,
				Slug:  val.Category.Slug,
			},
			User: entity.UserEntity{
				ID:   val.User.ID,
				Name: val.User.Name,
			},
			CreatedAt: val.CreatedAt,
		}

		resps = append(resps, resp)
	}

	return resps, nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
