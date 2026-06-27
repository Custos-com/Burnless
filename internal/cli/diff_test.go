package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffPreview_ValidFile(t *testing.T) {
	tmp := t.TempDir()
	sre := `
service: payments-api
team: platform-engineering
slos:
  - name: availability
    target: 99.9
    window: 30d
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
      - kubectl scale deploy/payments --replicas=+2
oncall:
  provider: pagerduty
  escalation_minutes: 10
  notify_slack: "#sre-incidents"
dashboards:
  provider: grafana
  auto_generate: true
  panels:
    - error_budget_remaining
`
	path := filepath.Join(tmp, "sre.yaml")
	if err := os.WriteFile(path, []byte(sre), 0o644); err != nil {
		t.Fatalf("failed to write sre.yaml: %v", err)
	}

	cmd := NewDiffCmd()
	cmd.SetArgs([]string{path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDiffTwoFiles_DetectsChanges(t *testing.T) {
	tmp := t.TempDir()

	old := `
service: payments-api
slos:
  - name: availability
    target: 99.9
    window: 30d
`
	new := `
service: payments-api
slos:
  - name: availability
    target: 99.95
    window: 30d
  - name: latency_p99
    target: 99.5
    window: 30d
`
	oldPath := filepath.Join(tmp, "old.yaml")
	newPath := filepath.Join(tmp, "new.yaml")

	if err := os.WriteFile(oldPath, []byte(old), 0o644); err != nil {
		t.Fatalf("failed to write old.yaml: %v", err)
	}
	if err := os.WriteFile(newPath, []byte(new), 0o644); err != nil {
		t.Fatalf("failed to write new.yaml: %v", err)
	}

	cmd := NewDiffCmd()
	cmd.SetArgs([]string{oldPath, newPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDiffTwoFiles_NoChanges(t *testing.T) {
	tmp := t.TempDir()
	sre := `
service: payments-api
slos:
  - name: availability
    target: 99.9
    window: 30d
`
	path1 := filepath.Join(tmp, "a.yaml")
	path2 := filepath.Join(tmp, "b.yaml")

	if err := os.WriteFile(path1, []byte(sre), 0o644); err != nil {
		t.Fatalf("failed to write a.yaml: %v", err)
	}
	if err := os.WriteFile(path2, []byte(sre), 0o644); err != nil {
		t.Fatalf("failed to write b.yaml: %v", err)
	}

	cmd := NewDiffCmd()
	cmd.SetArgs([]string{path1, path2})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
