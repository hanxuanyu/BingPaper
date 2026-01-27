package s3

import (
	"context"
	"fmt"
	"io"
	"strings"

	"BingPaper/internal/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client          *s3.Client
	bucket          string
	publicURLPrefix string
}

func NewS3Storage(endpoint, region, bucket, accessKey, secretKey, publicURLPrefix string, forcePathStyle bool) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = forcePathStyle
	})

	return &S3Storage{
		client:          client,
		bucket:          bucket,
		publicURLPrefix: publicURLPrefix,
	}, nil
}

func (s *S3Storage) Put(ctx context.Context, key string, r io.Reader, contentType string) (storage.StoredObject, error) {
	uploader := manager.NewUploader(s.client)
	output, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return storage.StoredObject{}, err
	}

	publicURL := ""
	if s.publicURLPrefix != "" {
		publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(s.publicURLPrefix, "/"), key)
	} else {
		publicURL = output.Location
	}

	return storage.StoredObject{
		Key:         key,
		ContentType: contentType,
		PublicURL:   publicURL,
	}, nil
}

func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", err
	}
	contentType := ""
	if output.ContentType != nil {
		contentType = *output.ContentType
	}
	return output.Body, contentType, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) PublicURL(key string) (string, bool) {
	if s.publicURLPrefix != "" {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(s.publicURLPrefix, "/"), key), true
	}
	// 也可以生成签名 URL，但这里简单处理
	return "", false
}
