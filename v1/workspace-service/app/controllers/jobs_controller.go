package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"workspace-service/app/models/dto"
	"workspace-service/app/services"
	"workspace-service/app/utils"
)

type JobsController struct {
	jobService services.JobService
	logger     *zap.Logger
}

func (wc *JobsController) CreateWorkspace(c *gin.Context) {
	body := dto.CreateJobRequest{
		ExecutorImage: "python",
	}

	if err := c.BindJSON(&body); err != nil {
		wc.logger.Error("Failed to bind json", zap.Error(err))
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}

	if err := utils.ValidateResourceId("projectId", body.ProjectId); err != nil {
		wc.logger.Warn("Rejected invalid job request", zap.Error(err))
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad Request"})
		return
	}
	if err := utils.ValidateExecutorImage(body.ExecutorImage); err != nil {
		wc.logger.Warn("Rejected invalid job request", zap.Error(err))
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad Request"})
		return
	}

	jobDetails, err := wc.jobService.CreateJob(body)
	if err != nil {
		wc.logger.Error("Failed to create job", zap.Error(err))
		c.AbortWithStatusJSON(
			500,
			gin.H{"error": "Internal Server Error"},
		)
		return
	}

	c.JSON(
		200,
		gin.H{"message": "success", "job": jobDetails},
	)
}

func NewJobsController(
	logger *zap.Logger,
	jobsService services.JobService,
) *JobsController {
	return &JobsController{
		jobService: jobsService,
		logger:     logger,
	}
}
