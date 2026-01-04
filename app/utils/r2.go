package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dariubs/scaffold/app/config"
	"github.com/google/uuid"
)

type R2Service struct {
	client *s3.Client
	bucket string
	region string
}

func NewR2Service() (*R2Service, error) {
	// Check if R2 configuration is available
	if config.C.CloudflareR2.AccountID == "" || config.C.CloudflareR2.AccessKeyID == "" ||
		config.C.CloudflareR2.SecretAccessKey == "" || config.C.CloudflareR2.Bucket == "" {
		return nil, fmt.Errorf("R2 configuration is incomplete")
	}

	accountID := config.C.CloudflareR2.AccountID
	accessKeyID := config.C.CloudflareR2.AccessKeyID
	secretAccessKey := config.C.CloudflareR2.SecretAccessKey
	bucket := config.C.CloudflareR2.Bucket
	region := config.C.CloudflareR2.Region

	// Create custom endpoint for Cloudflare R2
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	// Configure AWS SDK for R2
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &R2Service{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

// UploadFile uploads a file to R2 and returns the public URL
func (r2 *R2Service) UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s-%s%s", folder, uuid.New().String(), ext)
	key := fmt.Sprintf("%s/%s", folder, filename)

	// Upload to R2
	_, err = r2.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r2.bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
		ACL:         "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to R2: %v", err)
	}

	// Return public URL
	publicURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", r2.bucket, key)
	return publicURL, nil
}

// UploadProfileImage uploads a profile image and returns the URL
func (r2 *R2Service) UploadProfileImage(file *multipart.FileHeader, userID uint) (string, error) {
	return r2.UploadFile(file, "profiles")
}

// UploadImage uploads a general image and returns the URL
func (r2 *R2Service) UploadImage(file *multipart.FileHeader, folder string) (string, error) {
	return r2.UploadFile(file, folder)
}

// DeleteFile deletes a file from R2
func (r2 *R2Service) DeleteFile(fileURL string) error {
	// Extract key from URL
	// URL format: https://bucket.r2.cloudflarestorage.com/folder/filename
	parts := strings.Split(fileURL, "/")
	if len(parts) < 4 {
		return fmt.Errorf("invalid file URL format")
	}

	key := strings.Join(parts[3:], "/")

	_, err := r2.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r2.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from R2: %v", err)
	}

	return nil
}

// DeleteFileByKey deletes a file by its key
func (r2 *R2Service) DeleteFileByKey(key string) error {
	_, err := r2.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r2.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from R2: %v", err)
	}

	return nil
}

// GetFileURL returns the public URL for a file
func (r2 *R2Service) GetFileURL(key string) string {
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", r2.bucket, key)
}

// ListFiles lists files in a folder
func (r2 *R2Service) ListFiles(folder string) ([]string, error) {
	var files []string

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(r2.bucket),
		Prefix: aws.String(folder + "/"),
	}

	result, err := r2.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	for _, object := range result.Contents {
		if object.Key != nil {
			files = append(files, *object.Key)
		}
	}

	return files, nil
}

// FileExists checks if a file exists in R2
func (r2 *R2Service) FileExists(key string) (bool, error) {
	_, err := r2.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(r2.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
