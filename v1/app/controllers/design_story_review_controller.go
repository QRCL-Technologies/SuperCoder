package controllers

import (
	"ai-developer/app/services"
	"ai-developer/app/types/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DesignStoryReviewController struct {
	designStoryService *services.DesignStoryReviewService
	storyService       *services.StoryService
	projectService     *services.ProjectService
}

func NewDesignStoryReviewController(
	designStoryService *services.DesignStoryReviewService,
	storyService *services.StoryService,
	projectService *services.ProjectService,
) *DesignStoryReviewController {
	return &DesignStoryReviewController{
		designStoryService: designStoryService,
		storyService:       storyService,
		projectService:     projectService,
	}
}

func (ctrl *DesignStoryReviewController) CreateCommentForDesignStory(c *gin.Context) {
	var createCommentRequest request.CreateDesignStoryCommentRequest
	if err := c.ShouldBindJSON(&createCommentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, ok := authorizeStoryAccess(c, ctrl.storyService, ctrl.projectService, int(createCommentRequest.StoryID)); !ok {
		return
	}
	err := ctrl.designStoryService.CreateComment(createCommentRequest.StoryID, createCommentRequest.Comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
	return
}
