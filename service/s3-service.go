package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type IS3Service interface {
	Upload(filePath string) error
}

type S3Service struct {
	bucketName string
	client     *s3.Client
}

func NewS3Service(region string, bucketName string, accessKey string, secretKey string) *S3Service {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	return &S3Service{
		bucketName: bucketName,
		client:     s3Client,
	}
}

func (s S3Service) Upload(fileName string) error {
	file, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatalf("failed to resolve file path: %v", err)
	}

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   strings.NewReader(string(fileBytes)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}
