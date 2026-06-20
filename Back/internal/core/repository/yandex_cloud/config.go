package core_yandex_cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	KeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	Region    string `envconfig:"REGION"`
	Endpoint  string `envconfig:"ENDPOINT_URL"`
}

func newConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() aws.Config {
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get yandex_cloud client config: %w", err)
		panic(err)
	}

	cfgS3, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithBaseEndpoint(cfg.Endpoint),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.KeyID,
				cfg.SecretKey,
				"",
			),
		),
	)

	return cfgS3
}
