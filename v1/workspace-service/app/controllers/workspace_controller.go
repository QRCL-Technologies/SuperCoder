package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"workspace-service/app/config"
	"workspace-service/app/models/dto"
	"workspace-service/app/services"
	"workspace-service/app/utils"
)

type WorkspaceController struct {
	wsService services.WorkspaceService
	wsConfig  *config.WorkspaceServiceConfig
	logger    *zap.Logger
}

func (wc *WorkspaceController) rejectInvalid(c *gin.Context, err error) {
	wc.logger.Warn("Rejected invalid workspace request", zap.Error(err))
	c.AbortWithStatusJSON(400, gin.H{
		"error": "Bad Request",
	})
}

func (wc *WorkspaceController) validateCreateWorkspace(body *dto.CreateWorkspace) error {
	if err := utils.ValidateResourceId("workspaceId", body.WorkspaceId); err != nil {
		return err
	}
	if body.BackendTemplate == nil {
		return errors.New("backendTemplate is required")
	}
	if err := utils.ValidateTemplateName(utils.TemplatesRoot, *body.BackendTemplate); err != nil {
		return err
	}
	if body.FrontendTemplate != nil {
		if err := utils.ValidateTemplateName(utils.TemplatesRoot, *body.FrontendTemplate); err != nil {
			return err
		}
	}
	return utils.ValidateRemoteURL(body.RemoteURL, wc.wsConfig.GitnessHost())
}

func (wc *WorkspaceController) CreateWorkspace(c *gin.Context) {
	body := dto.CreateWorkspace{}
	if err := c.BindJSON(&body); err != nil {
		wc.logger.Error("Failed to bind json", zap.Error(err))
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
	if body.FrontendTemplate != nil && *body.FrontendTemplate == "" {
		body.FrontendTemplate = nil
	}
	if err := wc.validateCreateWorkspace(&body); err != nil {
		wc.rejectInvalid(c, err)
		return
	}
	wsDetails, err := wc.wsService.CreateWorkspace(body.WorkspaceId, *body.BackendTemplate, body.FrontendTemplate, body.RemoteURL, body.GitnessUserName, body.GitnessToken)
	if err != nil {
		c.AbortWithStatusJSON(
			500,
			gin.H{"error": "Internal Server Error"},
		)
		return
	}
	c.JSON(
		200,
		gin.H{"message": "success", "workspace": wsDetails},
	)
}

func (wc *WorkspaceController) CreateFrontendWorkspace(c *gin.Context) {
	nextjs := "nextjs"
	body := dto.CreateWorkspace{
		FrontendTemplate: &nextjs,
	}
	if err := c.BindJSON(&body); err != nil {
		wc.logger.Error("Failed to bind json", zap.Error(err))
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Bad Request",
		})
		return
	}
	if err := utils.ValidateResourceId("workspaceId", body.WorkspaceId); err != nil {
		wc.rejectInvalid(c, err)
		return
	}
	if err := utils.ValidateResourceId("storyHashId", body.StoryHashId); err != nil {
		wc.rejectInvalid(c, err)
		return
	}
	if body.FrontendTemplate == nil {
		wc.rejectInvalid(c, errors.New("frontendTemplate is required"))
		return
	}
	if err := utils.ValidateTemplateName(utils.TemplatesRoot, *body.FrontendTemplate); err != nil {
		wc.rejectInvalid(c, err)
		return
	}
	wsDetails, err := wc.wsService.CreateFrontendWorkspace(body.StoryHashId, body.WorkspaceId, *body.FrontendTemplate)
	if err != nil {
		c.AbortWithStatusJSON(
			500,
			gin.H{"error": "Internal Server Error"},
		)
		return
	}
	c.JSON(
		200,
		gin.H{"message": "success", "workspace": wsDetails},
	)
}

func (wc *WorkspaceController) DeleteWorkspace(c *gin.Context) {
	workspaceId := c.Param("workspaceId")
	if err := utils.ValidateResourceId("workspaceId", workspaceId); err != nil {
		wc.rejectInvalid(c, err)
		return
	}
	err := wc.wsService.DeleteWorkspace(workspaceId)
	if err != nil {
		c.AbortWithStatusJSON(
			500,
			gin.H{"error": "Internal Server Error"},
		)
		return
	}
	c.JSON(
		200,
		gin.H{"message": "success"},
	)

}

func NewWorkspaceController(
	logger *zap.Logger,
	wsService services.WorkspaceService,
	wsConfig *config.WorkspaceServiceConfig,
) *WorkspaceController {
	return &WorkspaceController{
		wsService: wsService,
		wsConfig:  wsConfig,
		logger:    logger,
	}
}
