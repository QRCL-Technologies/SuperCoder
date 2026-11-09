package controllers

import (
	"ai-developer/app/models"
	"ai-developer/app/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func currentUser(c *gin.Context) (*models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	user, ok := userInterface.(*models.User)
	return user, ok && user != nil
}

// authorizeProjectAccess aborts the request and returns false unless the
// project exists and belongs to the caller's organisation.
func authorizeProjectAccess(c *gin.Context, projectService *services.ProjectService, projectID int) (*models.Project, bool) {
	user, ok := currentUser(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil, false
	}
	project, err := projectService.GetProjectDetailsById(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return nil, false
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return nil, false
	}
	if project == nil || project.OrganisationID != user.OrganisationID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return nil, false
	}
	return project, true
}

// authorizeStoryAccess aborts the request and returns false unless the story
// exists and its project belongs to the caller's organisation.
func authorizeStoryAccess(c *gin.Context, storyService *services.StoryService, projectService *services.ProjectService, storyID int) (*models.Story, bool) {
	if _, ok := currentUser(c); !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil, false
	}
	story, err := storyService.GetStoryById(int64(storyID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return nil, false
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch story"})
		return nil, false
	}
	if story == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return nil, false
	}
	if _, ok := authorizeProjectAccess(c, projectService, int(story.ProjectID)); !ok {
		return nil, false
	}
	return story, true
}
