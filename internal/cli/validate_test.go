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

func TestRunValidate_DefaultFile(t *testing.T) {
	// change to a temp dir with a valid sre.yaml
	tmp := t.TempDir()
	sre := `
service: payments-api
slos:
  - name: availability
    target: 99.9
    window: 30d
`
	if err := os.WriteFile(filepath.Join(tmp, "sre.yaml"), []byte(sre), 0o644); err != nil {
		t.Fatalf("failed to write sre.yaml: %v", err)
	}

	// change working directory to temp dir
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	// run validate with NO arguments — should find sre.yaml automatically
	cmd := NewValidateCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error with default sre.yaml, got: %v", err)
	}
}
