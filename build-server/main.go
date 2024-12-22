package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

// runCommand executes a shell command and captures its output.
func runCommand(cmd string) error {
	fmt.Printf("Executing: %s\n", cmd)

	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", out.String())
		return err
	}
	fmt.Printf("Output: %s\n", out.String())
	return nil
}

// uploadFile uploads a single file to S3.
func uploadFile(ctx context.Context, client *s3.Client, bucket, key, filePath string) error {
	fmt.Printf("Uploading %s to S3 as %s...\n", filePath, key)

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
	fmt.Printf("Successfully uploaded %s\n", filePath)
	return nil
}

// uploadDirectory recursively uploads all files in a directory to S3.
func uploadDirectory(ctx context.Context, client *s3.Client, bucket, baseKey, dirPath string) error {
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
		if err := uploadFile(ctx, client, bucket, s3Key, path); err != nil {
			log.Printf("Failed to upload %s; continuing with next file.\n", path)
		}
		return nil
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" {
		log.Fatalf("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set in the environment.")
	}

	projectID := flag.String("pid", "", "Project ID to use for S3 key")
	repoLink := flag.String("repo", "", "Repository link")
	bucket := "solace-outputs"
	flag.Parse()

	if *projectID == "" {
		log.Fatalf("Project ID is required. Use the --pid flag to specify it.")
	}

	if *repoLink == "" {
		log.Fatalf("Repository link is required. Use the --repo flag to specify it.")
	}

	fmt.Println("Build Started...")

	if err := runCommand(fmt.Sprintf("git clone %s ./output", *repoLink)); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	fmt.Println("Repository cloned.")

	outDirPath := filepath.Join(".", "output")
	if err := runCommand(fmt.Sprintf("cd %s && npm install && npm run build", outDirPath)); err != nil {
		log.Fatalf("Failed to build the project: %v", err)
	}

	// Set up AWS S3 client
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			"",
		)),
		config.WithRegion("ap-south-1"),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}
	client := s3.NewFromConfig(cfg)
	fmt.Println("S3 client initialized.")

	distFolderPath := filepath.Join(outDirPath, "dist")
	baseKey := fmt.Sprintf("__outputs/%s", *projectID)

	if err := uploadDirectory(context.Background(), client, bucket, baseKey, distFolderPath); err != nil {
		log.Fatalf("Failed to upload files: %v", err)
	}

	fmt.Println("Build process completed successfully.")
}
