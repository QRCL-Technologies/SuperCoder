package controllers

import (
	"ai-developer/app/services"
	"ai-developer/app/types/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PullRequestCommentsController struct {
	pullRequestCommentService *services.PullRequestCommentsService
}

func NewPullRequestCommentController(pullRequestCommentService *services.PullRequestCommentsService) *PullRequestCommentsController {
	return &PullRequestCommentsController{
		pullRequestCommentService: pullRequestCommentService,
	}
}

func (ctrl *PullRequestCommentsController) CreateCommentForPrID(c *gin.Context) {
	pullRequestID, err := strconv.ParseUint(c.Param("pull_request_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pull request ID"})
		return
	}
	var createCommentRequest request.CreateCommentRequest
	if err := c.ShouldBindJSON(&createCommentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if createCommentRequest.PullRequestID != 0 && createCommentRequest.PullRequestID != uint(pullRequestID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pull request ID in body does not match path"})
		return
	}
	err = ctrl.pullRequestCommentService.CreateComment(uint(pullRequestID), createCommentRequest.Comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
	return
}
