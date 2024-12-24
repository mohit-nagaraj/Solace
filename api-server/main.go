package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type JobRequest struct {
	PID      string `json:"pid" binding:"required"`
	RepoLink string `json:"repo_link" binding:"required"`
}

type Config struct {
	ProjectID string
	Location  string
	Jobname   string
	ImageName string
}

func main() {
	config := Config{
		ProjectID: "artful-talon-445419-p7",
		Location:  "asia-south1",
		Jobname:   "build-image",
		ImageName: "asia-south1-docker.pkg.dev/artful-talon-445419-p7/solace/build-image:latest",
	}

	r := gin.Default()
	port := "9000"
	r.Run(":" + port)
}
