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
func runCommand(ctx context.Context, cmd string) error {
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
func uploadFile(ctx context.Context, client *s3.Client, bucket string, key string, filepath string) error {
	fmt.Printf("Uploading %s to S3...\n", filepath)

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", filepath, err)
		return err
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(mime.TypeByExtension(filepath)),
	})

	if err != nil {
		log.Printf("Failed to upload file %s: %v\n", filepath, err)
		return err
	}
	fmt.Printf("Successfully uploaded %s\n", filepath)
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" {
		log.Fatalf("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are required.")
		return
	}

	projectID := flag.String("pid", "", "Project ID to use for S3 key")
	bucket := "solace-outputs"
	flag.Parse()

	if *projectID == "" {
		log.Fatalf("Project ID is required.")
		return
	}

	fmt.Println("Build Started...")
	repoLink := "https://github.com/mohit-nagaraj/qr-generator"

	if err := runCommand(ctx, fmt.Sprintf("git clone %s ./output", repoLink)); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
		return
	}
	fmt.Println("Repository cloned.")

	outDirPath := filepath.Join(".", "output")
	if err := runCommand(ctx, fmt.Sprintf("cd %s", outDirPath)); err != nil {
		log.Fatalf("Build failed: %v", err)
		return
	}

	if err := runCommand(ctx, "npm install && npm run build"); err != nil {
		log.Fatalf("Build 2 failed: %v", err)
		return
	}

	distFolderPath := filepath.Join(outDirPath, "dist")
	files, err := os.ReadDir(distFolderPath)
	if err != nil {
		log.Fatalf("Failed to read dist folder: %v", err)
		return
	}

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
		log.Fatalf("Failed to load AWS config: %v\n", err)
		return
	}
	client := s3.NewFromConfig(cfg)
	fmt.Println("Initialized S3 client.")

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(distFolderPath, file.Name())
		s3Key := fmt.Sprintf("__outputs/%s/%s", *projectID, file.Name())

		if err := uploadFile(ctx, client, bucket, s3Key, filePath); err != nil {
			log.Println("Upload failed; continuing with next file.")
		}
	}

	fmt.Println("Build Completed.")
}
