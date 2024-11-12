package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/gofiber/fiber/v2/log"
)

func (cfg Config) LoadAwsConfig() aws.Config {
	conf, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.CloudflareR2.ApiKey, cfg.CloudflareR2.ApiSecret, "")),
		awsConfig.WithRegion("auto"))
	if err != nil {
		log.Fatalf("unable to load SDK AWS Config, %v", err)
	}

	log.Info("Success Load SDK AWS Config")

	return conf
}
