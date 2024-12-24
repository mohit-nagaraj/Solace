package main

import (
	"context"
	"fmt"
	"net/http"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
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

	r.POST("/run-job-override", handleRunJobOverride(config))

	port := "9000"
	r.Run(":" + port)
}

func handleRunJobOverride(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req JobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		client, err := run.NewJobsClient(ctx, option.WithCredentialsFile("gcp.json"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create client: %v", err)})
			return
		}
		defer client.Close()

		// Construct container arguments
		containerArgs := []string{
			fmt.Sprintf("--pid=%s", req.PID),
			fmt.Sprintf("--repo=%s", req.RepoLink),
		}

		// Create the job request with overrides
		jobReq := &runpb.RunJobRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s",
				config.ProjectID,
				config.Location,
				config.Jobname,
			),
			Overrides: &runpb.RunJobRequest_Overrides{
				ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
					{
						Name: "build-image",
						Args: containerArgs,
					},
				},
			},
		}

		operation, err := client.RunJob(ctx, jobReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to run job: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Job started successfully with overrides",
			"name":    operation.Name(),
		})
	}
}
