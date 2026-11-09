package config

import "github.com/knadh/koanf/v2"

type WorkspaceServiceConfig struct {
	config *koanf.Koanf
}

func (wsc *WorkspaceServiceConfig) GetEndpoint() string {
	return config.String("workspace.service.endpoint")
}

func (wsc *WorkspaceServiceConfig) GetAPIKey() string {
	return config.String("workspace.service.api.key")
}

func NewWorkspaceServiceConfig(config *koanf.Koanf) *WorkspaceServiceConfig {
	return &WorkspaceServiceConfig{config}
}
