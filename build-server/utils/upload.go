package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mohit-nagaraj/solace/build-server/logs"
)

// uploadFile uploads a single file to S3.
func uploadFile(ctx context.Context, projectID *string, client *s3.Client, bucket, key, filePath string) error {

	logs.PublishLog("logs:"+*projectID, fmt.Sprintf("Uploading %s to S3 as %s...\n", filePath, key))

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", filePath, err)
		return err
	}

	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		log.Printf("Failed to upload file %s: %v\n", filePath, err)
		return err
	}

	logs.PublishLog("logs:"+*projectID, fmt.Sprintf("Successfully uploaded %s\n", filePath))
	return nil
}

// uploadDirectory recursively uploads all files in a directory to S3.
func UploadDirectory(ctx context.Context, projectID *string, client *s3.Client, bucket, baseKey, dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %v", path, err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Compute relative path for S3 key
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path for %s: %v", path, err)
		}
		s3Key := filepath.Join(baseKey, relPath)

		// Upload the file
		if err := uploadFile(ctx, projectID, client, bucket, s3Key, path); err != nil {
			log.Printf("Failed to upload %s; continuing with next file.\n", path)
		}
		return nil
	})
}
