package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runCommand(ctx context.Context, cmd string) error {
	fmt.Printf("Executing: %s", cmd)

	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Error: %s", out.String())
		return err
	}
	fmt.Printf("Output: %s", out.String())
	return nil
}

func main() {
	ctx := context.Background()

	fmt.Print("Build Started...\n")

	fmt.Print("Cloning repository...\n")
	repoLink := "https://github.com/your/repo.git"
	if err := runCommand(ctx, fmt.Sprintf("git clone %s ./output", repoLink)); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	fmt.Print("Repository cloned\n")

	outDirPath := filepath.Join(".", "output")
	if err := runCommand(ctx, fmt.Sprintf("cd %s && npm install && npm run build", outDirPath)); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	distFolderPath := filepath.Join(outDirPath, "dist")
	files, err := os.ReadDir(distFolderPath)
	if err != nil {
		log.Fatalf("Failed to read dist folder: %v", err)
	}

	fmt.Print("Uploading files to S3...\n")

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(distFolderPath, file.Name())

		fmt.Printf("Uploading %s to S3...\n", filePath)
	}

	fmt.Print("Build Completed\n")
}
