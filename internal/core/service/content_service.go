package service

import (
	"context"
	"latihan-portal-news/config"
	"latihan-portal-news/internal/adapter/cloudflare"
	"latihan-portal-news/internal/adapter/repository"
	"latihan-portal-news/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context, query *entity.QueryString) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContentByID(ctx context.Context, id int64) error

	GetContentByCategoryID(ctx context.Context, categoryID int64, search string) ([]entity.ContentEntity, error)
	UploadImageCloudFlareR2(ctx context.Context, req entity.FileUploadEntity) (string, error)
}

type contentService struct {
	contentRepository repository.ContentRepository
	cfg               *config.Config
	cloudflareR2      cloudflare.FileRepository
}

// UploadImageCloudFlareR2 implements ContentService.
func (c *contentService) UploadImageCloudFlareR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	urlImage, err := c.cloudflareR2.UploadFile(&req)
	if err != nil {
		code = "[Service] UploadImageCloudFlareR2 - 1"
		log.Errorw(code, err)
		return "", err
	}
	return urlImage, nil
}

// GetContentByCategoryID implements ContentService.
func (c *contentService) GetContentByCategoryID(ctx context.Context, categoryID int64, search string) ([]entity.ContentEntity, error) {
	results, err := c.contentRepository.GetContentByCategoryID(ctx, categoryID, search)
	if err != nil {
		code = "[Service] GetContentByCategoryID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return results, nil
}

// DeleteContentByID implements ContentService.
func (c *contentService) DeleteContentByID(ctx context.Context, id int64) error {
	err = c.contentRepository.DeleteContentByID(ctx, id)
	if err != nil {
		code = "[Service] DeleteContentByID - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// UpdateContent implements ContentService.
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	result, err := c.contentRepository.GetContentByID(ctx, req.ID)
	if err != nil {
		code = "[Service] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	if req.Image == "" {
		req.Image = result.Image
	}

	err = c.contentRepository.UpdateContent(ctx, req)
	if err != nil {
		code = "[Service] UpdateContent - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// CreateContent implements ContentService.
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepository.CreateContent(ctx, req)
	if err != nil {
		code = "[Service] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentService.
func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	result, err := c.contentRepository.GetContentByID(ctx, id)
	if err != nil {
		code = "[Service] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

// GetContents implements ContentService.
func (c *contentService) GetContents(ctx context.Context, query *entity.QueryString) ([]entity.ContentEntity, error) {
	results, err := c.contentRepository.GetContents(ctx, query)
	if err != nil {
		code = "[Service] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return results, nil
}

func NewContentService(contentRepository repository.ContentRepository, cfg *config.Config, clodFlareRepo cloudflare.FileRepository) ContentService {
	return &contentService{
		contentRepository: contentRepository,
		cfg:               cfg,
		cloudflareR2:      clodFlareRepo,
	}
}
