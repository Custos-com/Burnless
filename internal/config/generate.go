package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GenerateDefault creates a new SREConfig with sensible defaults
// for the given service name and team.
func GenerateDefault(serviceName string, team string) *SREConfig {
	return &SREConfig{
		Service: serviceName,
		Team:    team,

		// One SLO — availability at 99.9%
		SLOs: []SLO{
			{
				Name:   "availability",
				Target: 99.9,
				Window: "30d",
				Indicator: Indicator{
					Metric:     "http_requests_total",
					GoodFilter: `status!~"5.."`,
				},
			},
		},

		// Two burn rate alerts — critical and warning
		ErrorBudget: ErrorBudget{
			BurnRateAlerts: []BurnRateAlert{
				{
					Severity:  "critical",
					Rate:      14.4,
					Window:    "1h",
					Remediate: "scale-up",
				},
				{
					Severity:  "warning",
					Rate:      6.0,
					Window:    "6h",
					Remediate: "restart-pods",
				},
			},
		},

		// Two runbooks — one auto, one semi-auto
		Runbooks: map[string]Runbook{
			"scale-up": {
				Mode: "auto",
				Steps: []string{
					fmt.Sprintf("kubectl scale deploy/%s --replicas=+2", serviceName),
					"wait 60s",
					"assert slo.availability > 99.5",
				},
			},
			"restart-pods": {
				Mode: "semi-auto",
				Steps: []string{
					fmt.Sprintf("kubectl rollout restart deploy/%s", serviceName),
				},
			},
		},

		// On-call config
		OnCall: OnCallConfig{
			Provider:          "pagerduty",
			EscalationMinutes: 10,
			NotifySlack:       "#sre-incidents",
		},

		// Dashboard config
		Dashboards: DashboardConfig{
			Provider:     "grafana",
			AutoGenerate: true,
			Panels: []string{
				"error_budget_remaining",
				"burn_rate_1h",
				"burn_rate_6h",
				"request_success_rate",
			},
		},
	}
}

// WriteToFile serialises cfg to YAML and writes it to path.
func WriteToFile(cfg *SREConfig, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("serialising config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
