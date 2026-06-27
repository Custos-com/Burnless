package cli

import (
	"fmt"

	"github.com/burnless/burnless/internal/config"

	"github.com/spf13/cobra"
)

// NewDiffCmd creates the "burnless diff" command.
func NewDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "diff [file] [new-file]",
		Short:        "Show what burnless apply would create, or diff two sre.yaml files",
		Args:         cobra.MaximumNArgs(2),
		SilenceUsage: true,
		RunE:         runDiff,
	}
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Option A — two files provided: compare them
	if len(args) == 2 {
		return runFileDiff(args[0], args[1])
	}

	// Option C — zero or one file: show what apply would do
	path := "sre.yaml"
	if len(args) == 1 {
		path = args[0]
	}
	return runApplyPreview(path)
}

// runApplyPreview shows what "burnless apply" would create from sre.yaml.
func runApplyPreview(path string) error {
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}

	fmt.Printf("Loaded: %s\n", path)
	fmt.Printf("Service: %s\n", cfg.Service)
	if cfg.Team != "" {
		fmt.Printf("Team:    %s\n", cfg.Team)
	}
	fmt.Println()
	fmt.Println("Would create:")

	// SLOs
	for _, slo := range cfg.SLOs {
		fmt.Printf("  + SLO: %s (target: %.1f%%, window: %s)\n",
			slo.Name, slo.Target, slo.Window)
	}

	// Prometheus alerting rules
	alertCount := len(cfg.ErrorBudget.BurnRateAlerts)
	if alertCount > 0 {
		fmt.Printf("  + %d Prometheus alerting rule(s)\n", alertCount)
	}

	// Runbooks
	autoCount := 0
	semiCount := 0
	for _, rb := range cfg.Runbooks {
		if rb.Mode == "auto" {
			autoCount++
		} else {
			semiCount++
		}
	}
	if len(cfg.Runbooks) > 0 {
		fmt.Printf("  + %d runbook(s) (%d auto, %d semi-auto)\n",
			len(cfg.Runbooks), autoCount, semiCount)
	}

	// Grafana dashboard
	if cfg.Dashboards.AutoGenerate {
		fmt.Printf("  + 1 Grafana dashboard (%d panel(s)) via %s\n",
			len(cfg.Dashboards.Panels), cfg.Dashboards.Provider)
	}

	// On-call
	if cfg.OnCall.Provider != "" {
		fmt.Printf("  + 1 on-call policy via %s (escalation: %d min)\n",
			cfg.OnCall.Provider, cfg.OnCall.EscalationMinutes)
	}

	fmt.Println()
	fmt.Println("Run 'burnless apply' to deploy these changes.")
	return nil
}

// runFileDiff compares two sre.yaml files and shows what changed.
func runFileDiff(oldPath, newPath string) error {
	oldCfg, err := config.Load(oldPath)
	if err != nil {
		return fmt.Errorf("loading %s: %w", oldPath, err)
	}

	newCfg, err := config.Load(newPath)
	if err != nil {
		return fmt.Errorf("loading %s: %w", newPath, err)
	}

	fmt.Printf("Comparing %s → %s\n\n", oldPath, newPath)

	changed := false

	// service
	if oldCfg.Service != newCfg.Service {
		fmt.Printf("  ~ service: %s → %s\n", oldCfg.Service, newCfg.Service)
		changed = true
	} else {
		fmt.Printf("    service: %s (unchanged)\n", newCfg.Service)
	}

	// team
	if oldCfg.Team != newCfg.Team {
		fmt.Printf("  ~ team: %s → %s\n", oldCfg.Team, newCfg.Team)
		changed = true
	}

	// SLOs
	oldSLOs := map[string]config.SLO{}
	for _, s := range oldCfg.SLOs {
		oldSLOs[s.Name] = s
	}
	newSLOs := map[string]config.SLO{}
	for _, s := range newCfg.SLOs {
		newSLOs[s.Name] = s
	}

	for name, ns := range newSLOs {
		if os, exists := oldSLOs[name]; !exists {
			fmt.Printf("  + SLO: %s added (target: %.1f%%)\n", name, ns.Target)
			changed = true
		} else if os.Target != ns.Target {
			fmt.Printf("  ~ SLO: %s target: %.2f%% → %.2f%%\n", name, os.Target, ns.Target)
			changed = true
		} else {
			fmt.Printf("    SLO: %s (unchanged)\n", name)
		}
	}
	for name := range oldSLOs {
		if _, exists := newSLOs[name]; !exists {
			fmt.Printf("  - SLO: %s removed\n", name)
			changed = true
		}
	}

	// runbooks
	for name := range newCfg.Runbooks {
		if _, exists := oldCfg.Runbooks[name]; !exists {
			fmt.Printf("  + runbook: %s added\n", name)
			changed = true
		}
	}
	for name := range oldCfg.Runbooks {
		if _, exists := newCfg.Runbooks[name]; !exists {
			fmt.Printf("  - runbook: %s removed\n", name)
			changed = true
		}
	}

	fmt.Println()
	if !changed {
		fmt.Println("No changes detected.")
	} else {
		fmt.Println("Changes detected. Review before running 'burnless apply'.")
	}

	return nil
}
