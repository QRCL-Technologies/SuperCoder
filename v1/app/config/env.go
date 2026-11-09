package config

import (
	"ai-developer/app/constants"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type EnvConfig struct {
	config *koanf.Koanf
	logger *zap.Logger
}

func (e EnvConfig) IsDevelopment() bool {
	return e.config.String("app.env") == constants.Development
}

// DefaultUserLoginEnabled reports whether the passwordless default-user login is allowed.
// It requires both development mode and an explicit opt-in flag.
func (e EnvConfig) DefaultUserLoginEnabled() bool {
	return e.IsDevelopment() && e.config.Bool("auth.default.user.enabled")
}

func (e EnvConfig) Domain() string {
	var domain string
	if !e.IsDevelopment() {
		domain = e.config.String("domain")
	}
	e.logger.Debug("Setting domain for auth", zap.String("domain", domain))
	return domain
}

func NewEnvConfig(config *koanf.Koanf, logger *zap.Logger) *EnvConfig {
	return &EnvConfig{
		config: config,
		logger: logger.Named("EnvConfig"),
	}
}
