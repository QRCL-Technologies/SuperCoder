package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateResourceId(t *testing.T) {
	valid := []string{"abcde-6f1c2a9e-1b2c-4d5e-8f90-1234567890ab", "a", "x1y2z"}
	for _, id := range valid {
		if err := ValidateResourceId("workspaceId", id); err != nil {
			t.Errorf("expected %q to be valid: %v", id, err)
		}
	}
	invalid := []string{"", "..", "../x", "a/b", "a..b/c", "/etc", "-abc", "abc-", "ABC", "a b", "a\x00b"}
	for _, id := range invalid {
		if err := ValidateResourceId("workspaceId", id); err == nil {
			t.Errorf("expected %q to be invalid", id)
		}
	}
}

func TestValidateTemplateName(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "flask"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTemplateName(root, "flask"); err != nil {
		t.Errorf("expected flask to be valid: %v", err)
	}
	invalid := []string{"", "missing", "file", "..", "../workspaces/x", "flask/../..", "/etc", "fl.ask"}
	for _, name := range invalid {
		if err := ValidateTemplateName(root, name); err == nil {
			t.Errorf("expected %q to be invalid", name)
		}
	}
}

func TestValidateRemoteURL(t *testing.T) {
	host := "gitness:3000"
	valid := []string{
		"http://admin:token@gitness:3000/git/space/project.git",
		"https://gitness:3000/git/space/project.git",
	}
	for _, u := range valid {
		if err := ValidateRemoteURL(u, host); err != nil {
			t.Errorf("expected %q to be valid: %v", u, err)
		}
	}
	invalid := []string{
		"",
		"http://evil.example.com/repo.git",
		"http://gitness:3000@evil.example.com/repo.git",
		"file:///etc/passwd",
		"ssh://gitness:3000/repo.git",
		"http://gitness:3001/repo.git",
	}
	for _, u := range invalid {
		if err := ValidateRemoteURL(u, host); err == nil {
			t.Errorf("expected %q to be invalid", u)
		}
	}
	if err := ValidateRemoteURL(valid[0], ""); err == nil {
		t.Error("expected empty allowed host to reject")
	}
}
