package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDefault(t *testing.T) {
	cfg := GenerateDefault("payments-api", "platform-team")

	if cfg.Service != "payments-api" {
		t.Fatalf("Service = %q, want payments-api", cfg.Service)
	}
	if cfg.Team != "platform-team" {
		t.Fatalf("Team = %q, want platform-team", cfg.Team)
	}
	if len(cfg.SLOs) != 1 {
		t.Fatalf("expected 1 SLO, got %d", len(cfg.SLOs))
	}
	if cfg.SLOs[0].Target != 99.9 {
		t.Fatalf("SLO target = %v, want 99.9", cfg.SLOs[0].Target)
	}
	if len(cfg.Runbooks) != 2 {
		t.Fatalf("expected 2 runbooks, got %d", len(cfg.Runbooks))
	}

	// the generated config should pass its own validation
	if err := cfg.Validate(); err != nil {
		t.Fatalf("generated config failed validation: %v", err)
	}
}

func TestWriteToFile(t *testing.T) {
	cfg := GenerateDefault("checkout-api", "payments-team")

	path := filepath.Join(t.TempDir(), "sre.yaml")
	if err := WriteToFile(cfg, path); err != nil {
		t.Fatalf("WriteToFile() error: %v", err)
	}

	// prove the file exists and is not empty
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("written file is empty")
	}

	// prove the file can be loaded back correctly
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() of written file failed: %v", err)
	}
	if loaded.Service != "checkout-api" {
		t.Fatalf("loaded Service = %q, want checkout-api", loaded.Service)
	}
}
