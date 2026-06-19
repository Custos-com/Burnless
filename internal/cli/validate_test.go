package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunValidate_ValidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sre.yaml")
	content := `
service: checkout-api
slos:
  - name: availability
    target: 99.9
    window: 30d
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd := NewValidateCmd()
	cmd.SetArgs([]string{path})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error for valid file, got: %v", err)
	}
}

func TestRunValidate_InvalidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sre.yaml")
	content := `team: missing-service-name`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd := NewValidateCmd()
	cmd.SetArgs([]string{path})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid file, got nil")
	}
}

func TestRunValidate_MissingFile(t *testing.T) {
	cmd := NewValidateCmd()
	cmd.SetArgs([]string{"/tmp/this-file-does-not-exist.yaml"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
