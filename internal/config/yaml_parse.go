package config

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// parseFlexibleFloat parses YAML scalars like 99.9, "99.5%", or "14.4x".
func parseFlexibleFloat(node *yaml.Node, suffix byte) (float64, error) {
	if node == nil {
		return 0, fmt.Errorf("missing value")
	}

	switch node.Kind {
	case yaml.ScalarNode:
		var f float64
		if err := node.Decode(&f); err == nil {
			return f, nil
		}

		s := strings.TrimSpace(node.Value)
		if suffix != 0 {
			s = strings.TrimSuffix(s, string(suffix))
			s = strings.TrimSpace(s)
		}

		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %q", node.Value)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("expected scalar, got %v", node.Kind)
	}
}

func (s *SLO) UnmarshalYAML(node *yaml.Node) error {
	var raw struct {
		Name        string    `yaml:"name"`
		Target      yaml.Node `yaml:"target"`
		Window      string    `yaml:"window"`
		Indicator   Indicator `yaml:"indicator"`
		ThresholdMS int       `yaml:"threshold_ms,omitempty"`
	}
	if err := node.Decode(&raw); err != nil {
		return err
	}

	target, err := parseFlexibleFloat(&raw.Target, '%')
	if err != nil {
		return fmt.Errorf("target: %w", err)
	}

	*s = SLO{
		Name:        raw.Name,
		Target:      target,
		Window:      raw.Window,
		Indicator:   raw.Indicator,
		ThresholdMS: raw.ThresholdMS,
	}
	return nil
}

func (a *BurnRateAlert) UnmarshalYAML(node *yaml.Node) error {
	var raw struct {
		Severity  string    `yaml:"severity"`
		Rate      yaml.Node `yaml:"rate"`
		Window    string    `yaml:"window"`
		Remediate string    `yaml:"remediate"`
	}
	if err := node.Decode(&raw); err != nil {
		return err
	}

	rate, err := parseFlexibleFloat(&raw.Rate, 'x')
	if err != nil {
		return fmt.Errorf("rate: %w", err)
	}

	*a = BurnRateAlert{
		Severity:  raw.Severity,
		Rate:      rate,
		Window:    raw.Window,
		Remediate: raw.Remediate,
	}
	return nil
}
