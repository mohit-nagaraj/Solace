package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mohit-nagaraj/solace/build-server/logs"
	"github.com/mohit-nagaraj/solace/build-server/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

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

	logs.PublishLog("logs:"+*projectID, fmt.Sprintf("Build started for project %s", *projectID))
	logs.PublishLog("logs:"+*projectID, "Cloning repository: "+*repoLink)

	if err := utils.RunCommand(fmt.Sprintf("git clone %s ./output", *repoLink)); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	logs.PublishLog("logs:"+*projectID, fmt.Sprintf("Repository cloned..."))

	outDirPath := filepath.Join(".", "output")
	logs.PublishLog("logs:"+*projectID, "cd ./output && npm install && npm run build")
	if err := utils.RunCommand(fmt.Sprintf("cd %s && npm install && npm run build", outDirPath)); err != nil {
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
	logs.PublishLog("logs:"+*projectID, fmt.Sprintf("S3 client initialized."))

	distFolderPath := filepath.Join(outDirPath, "dist")
	baseKey := fmt.Sprintf("__outputs/%s", *projectID)

	if err := utils.UploadDirectory(context.Background(), projectID, client, bucket, baseKey, distFolderPath); err != nil {
		log.Fatalf("Failed to upload files: %v", err)
	}

	logs.PublishLog("logs:"+*projectID, fmt.Sprintf("Build process completed successfully."))
}
