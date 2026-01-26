package s3

import (
	"context"
	"fmt"
	"io"
	"strings"

	"BingDailyImage/internal/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	session         *session.Session
	client          *s3.S3
	bucket          string
	publicURLPrefix string
}

func NewS3Storage(endpoint, region, bucket, accessKey, secretKey, publicURLPrefix string, forcePathStyle bool) (*S3Storage, error) {
	config := &aws.Config{
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(forcePathStyle),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return &S3Storage{
		session:         sess,
		client:          s3.New(sess),
		bucket:          bucket,
		publicURLPrefix: publicURLPrefix,
	}, nil
}

func (s *S3Storage) Put(ctx context.Context, key string, r io.Reader, contentType string) (storage.StoredObject, error) {
	uploader := s3manager.NewUploader(s.session)
	output, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
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
	output, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", err
	}
	return output.Body, aws.StringValue(output.ContentType), nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
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
