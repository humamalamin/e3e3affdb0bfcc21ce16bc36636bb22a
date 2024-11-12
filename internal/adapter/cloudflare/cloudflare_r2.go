package cloudflare

import (
	"context"
	"fmt"
	"latihan-portal-news/config"
	"latihan-portal-news/internal/core/domain/entity"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2/log"
)

var code string
var err error

type FileRepository interface {
	UploadFile(req *entity.FileUploadEntity) (string, error)
}

type CloudFlareR2Repository struct {
	Client  *s3.Client
	Bucket  string
	BaseUrl string
}

// UploadFile implements FileRepository.
func (c *CloudFlareR2Repository) UploadFile(req *entity.FileUploadEntity) (string, error) {
	openedFile, err := os.Open(req.Path)
	if err != nil {
		code = "[Cloudflare] CloudFlareR2Repository - 1"
		log.Errorw(code, err)
		return "", err
	}

	defer openedFile.Close()

	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(req.Name),
		Body:        openedFile,
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		code = "[Cloudflare] CloudFlareR2Repository - 2"
		log.Errorw(code, err)
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s", c.BaseUrl, req.Name)

	return fileURL, nil
}

func NewCloudflareR2Repository(client *s3.Client, cfg *config.Config) FileRepository {
	clientBase := s3.NewFromConfig(cfg.LoadAwsConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.CloudflareR2.AccountID))
	})

	baseUrl := cfg.CloudflareR2.PublicUrl
	return &CloudFlareR2Repository{
		Client:  clientBase,
		Bucket:  cfg.CloudflareR2.BucketName,
		BaseUrl: baseUrl,
	}
}
