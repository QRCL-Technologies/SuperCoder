package config

import (
	"errors"
	"fmt"
	"github.com/knadh/koanf/v2"
	"log"
	"time"
)

const MinJWTSecretLength = 32

type JWTConfig struct {
	config *koanf.Koanf
}

func (j JWTConfig) Secret() string {
	return j.config.String("jwt.secret.key")
}

func (j JWTConfig) ExpiryHours() time.Duration {
	expiryHours := j.config.String("jwt.expiry.hours")
	expiryDuration, err := time.ParseDuration(expiryHours)
	if err != nil {
		log.Fatalf("could not parse jwt.expiry_hours: %v", err)
	}
	return expiryDuration
}

func ValidateJWTSecret(secret string) error {
	if secret == "" {
		return errors.New("jwt.secret.key is not set: set AI_DEVELOPER_JWT_SECRET_KEY to a random value of at least 32 bytes")
	}
	if len(secret) < MinJWTSecretLength {
		return fmt.Errorf("jwt.secret.key is too short: AI_DEVELOPER_JWT_SECRET_KEY must be at least %d bytes", MinJWTSecretLength)
	}
	return nil
}

func NewJWTConfig(config *koanf.Koanf) (*JWTConfig, error) {
	if err := ValidateJWTSecret(config.String("jwt.secret.key")); err != nil {
		return nil, err
	}
	return &JWTConfig{config: config}, nil
}

func JWTSecret() string { return config.String("jwt.secret.key") }

func JWTExpiryHours() time.Duration {
	expiryHours := config.String("jwt.expiry.hours")
	expiryDuration, err := time.ParseDuration(expiryHours)
	if err != nil {
		log.Fatalf("could not parse jwt.expiry_hours: %v", err)
	}
	return expiryDuration
}
