package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSREConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "sre.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return path
}

func TestLoadValidSREConfig(t *testing.T) {
	path := writeSREConfig(t, `
service: checkout-api
team: platform
slos:
  - name: availability
    target: 99.9
    window: 30d
    indicator:
      metric: http_requests_total
      good_filter: 'status!~"5.."'
error_budget:
  burn_rate_alerts:
    - severity: critical
      rate: 14.4
      window: 1h
      remediate: scale-up
runbooks:
  scale-up:
    mode: auto
    steps:
      - kubectl scale deploy/checkout --replicas=+2
oncall:
  provider: pagerduty
  escalation_minutes: 10
  notify_slack: "#sre-incidents"
dashboards:
  provider: grafana
  auto_generate: true
  panels:
    - error_budget_remaining
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.Service != "checkout-api" {
		t.Fatalf("service = %q, want checkout-api", cfg.Service)
	}
	if cfg.Team != "platform" {
		t.Fatalf("team = %q, want platform", cfg.Team)
	}
	if len(cfg.SLOs) != 1 {
		t.Fatalf("slos length = %d, want 1", len(cfg.SLOs))
	}
	slo := cfg.SLOs[0]
	if slo.Name != "availability" || slo.Target != 99.9 || slo.Window != "30d" {
		t.Fatalf("slo = %#v, want availability target 99.9 window 30d", slo)
	}
	if slo.Indicator.Metric != "http_requests_total" {
		t.Fatalf("indicator metric = %q, want http_requests_total", slo.Indicator.Metric)
	}
	if len(cfg.ErrorBudget.BurnRateAlerts) != 1 {
		t.Fatalf("burn rate alerts length = %d, want 1", len(cfg.ErrorBudget.BurnRateAlerts))
	}
	if cfg.Runbooks["scale-up"].Mode != "auto" {
		t.Fatalf("runbook mode = %q, want auto", cfg.Runbooks["scale-up"].Mode)
	}
	if cfg.OnCall.Provider != "pagerduty" {
		t.Fatalf("oncall provider = %q, want pagerduty", cfg.OnCall.Provider)
	}
	if !cfg.Dashboards.AutoGenerate {
		t.Fatal("dashboards auto_generate = false, want true")
	}
}

func TestLoadRejectsMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name: "missing service",
			content: `
slos:
  - name: availability
    target: 99.9
    window: 30d
`,
			wantErr: "service name is required",
		},
		{
			name: "missing slos",
			content: `
service: checkout-api
`,
			wantErr: "at least one SLO must be defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(writeSREConfig(t, tt.content))
			if err == nil {
				t.Fatal("Load() returned nil error, want validation error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Load() error = %q, want to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRejectsInvalidSLOTargets(t *testing.T) {
	tests := []struct {
		name   string
		target float64
	}{
		{name: "below zero", target: -0.1},
		{name: "above one hundred", target: 100.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := SREConfig{
				Service: "checkout-api",
				SLOs: []SLO{{
					Name:   "availability",
					Target: tt.target,
					Window: "30d",
				}},
			}
			err := cfg.Validate()
			if err == nil {
				t.Fatal("Validate() returned nil error, want target validation error")
			}
			if !strings.Contains(err.Error(), "target must be between 0 and 100") {
				t.Fatalf("Validate() error = %q, want target range error", err)
			}
		})
	}
}

func TestValidateAllowsBoundarySLOTargets(t *testing.T) {
	cfg := SREConfig{
		Service: "checkout-api",
		SLOs: []SLO{
			{Name: "zero", Target: 0, Window: "30d"},
			{Name: "one hundred", Target: 100, Window: "30d"},
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
}
