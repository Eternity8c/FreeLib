package book_s3_repository

import (
	"context"
	"fmt"
	"mime/multipart"

	core_yandex_cloud "github.com/Eternity8c/FreeLib/internal/core/repository/yandex_cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BookS3Repository struct {
	client *core_yandex_cloud.Client
}

func NewBookS3Repository(client *core_yandex_cloud.Client) *BookS3Repository {
	return &BookS3Repository{
		client: client,
	}
}

const (
	bucketName = "tes-freelib-server"
)

func (r *BookS3Repository) SaveBookFile(ctx context.Context, file multipart.File, fileName string) (string, error) {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String("application/epub+zip"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", bucketName, fileName)

	return fileURL, nil
}
