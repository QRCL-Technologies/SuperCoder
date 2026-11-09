package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const TemplatesRoot = "/templates"

var (
	resourceIdPattern   = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
	templateNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,62}$`)
)

// ValidateResourceId checks that an id (workspace, story or project hash id) is a single
// DNS-label-like path segment, so it cannot escape /workspaces or be used as a k8s name injection.
func ValidateResourceId(field string, id string) error {
	if !resourceIdPattern.MatchString(id) {
		return fmt.Errorf("invalid %s", field)
	}
	return nil
}

// ValidateTemplateName checks that a template name is a plain name of an existing directory directly under templatesRoot.
func ValidateTemplateName(templatesRoot string, name string) error {
	if !templateNamePattern.MatchString(name) {
		return fmt.Errorf("invalid template name")
	}
	info, err := os.Stat(filepath.Join(templatesRoot, name))
	if err != nil || !info.IsDir() {
		return fmt.Errorf("unknown template")
	}
	return nil
}

// ValidateExecutorImage checks that an executor image key is a plain name.
func ValidateExecutorImage(name string) error {
	if !templateNamePattern.MatchString(name) {
		return fmt.Errorf("invalid executor image")
	}
	return nil
}

// ValidateRemoteURL checks that remoteURL is an http(s) URL pointing at allowedHost.
func ValidateRemoteURL(remoteURL string, allowedHost string) error {
	if allowedHost == "" {
		return fmt.Errorf("no allowed git host configured")
	}
	u, err := url.Parse(remoteURL)
	if err != nil {
		return fmt.Errorf("invalid remote url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid remote url scheme")
	}
	if !strings.EqualFold(u.Host, allowedHost) {
		return fmt.Errorf("remote url host not allowed")
	}
	return nil
}
