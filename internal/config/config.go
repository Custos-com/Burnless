// Package config handles parsing and validating sre.yaml files.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and parses an sre.yaml file from the given path.
func Load(path string) (*SREConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var cfg SREConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating %s: %w", path, err)
	}

	return &cfg, nil
}

// SREConfig is the top-level structure of sre.yaml.
type SREConfig struct {
	Service     string        `yaml:"service"`
	Team        string        `yaml:"team"`
	SLOs        []SLO         `yaml:"slos"`
	ErrorBudget ErrorBudget   `yaml:"error_budget"`
	Runbooks    map[string]Runbook `yaml:"runbooks"`
	OnCall      OnCallConfig  `yaml:"oncall"`
	Dashboards  DashboardConfig `yaml:"dashboards"`
}

// SLO defines a Service Level Objective.
type SLO struct {
	Name        string    `yaml:"name"`
	Target      float64   `yaml:"target"`
	Window      string    `yaml:"window"`
	Indicator   Indicator `yaml:"indicator"`
	ThresholdMS int       `yaml:"threshold_ms,omitempty"`
}

// Indicator defines the metric source for an SLO.
type Indicator struct {
	Metric     string `yaml:"metric"`
	GoodFilter string `yaml:"good_filter"`
}

// ErrorBudget holds burn rate alert configuration.
type ErrorBudget struct {
	BurnRateAlerts []BurnRateAlert `yaml:"burn_rate_alerts"`
}

// BurnRateAlert defines a burn rate threshold and remediation action.
type BurnRateAlert struct {
	Severity   string  `yaml:"severity"`
	Rate       float64 `yaml:"rate"`
	Window     string  `yaml:"window"`
	Remediate  string  `yaml:"remediate"`
}

// Runbook defines an executable set of remediation steps.
type Runbook struct {
	Mode  string   `yaml:"mode"` // auto | semi-auto
	Steps []string `yaml:"steps"`
}

// OnCallConfig defines on-call routing.
type OnCallConfig struct {
	Provider           string `yaml:"provider"`
	EscalationMinutes  int    `yaml:"escalation_minutes"`
	NotifySlack        string `yaml:"notify_slack"`
}

// DashboardConfig defines dashboard generation settings.
type DashboardConfig struct {
	Provider     string   `yaml:"provider"`
	AutoGenerate bool     `yaml:"auto_generate"`
	Panels       []string `yaml:"panels"`
}

// Validate checks the SREConfig for required fields and valid values.
func (c *SREConfig) Validate() error {
	if c.Service == "" {
		return fmt.Errorf("service name is required")
	}
	if len(c.SLOs) == 0 {
		return fmt.Errorf("at least one SLO must be defined")
	}
	for i, slo := range c.SLOs {
		if slo.Target < 0 || slo.Target > 100 {
			return fmt.Errorf("slo[%d] target must be between 0 and 100, got %.2f", i, slo.Target)
		}
	}
	return nil
}
