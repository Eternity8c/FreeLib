package core_yandex_cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	*s3.Client
}

func NewClient(ctx context.Context, config aws.Config) (*Client, error) {
	client := s3.NewFromConfig(config)

	_, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("s3 ping failed: %w", err)
	}

	return &Client{
		Client: client,
	}, nil
}
