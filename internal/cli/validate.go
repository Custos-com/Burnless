package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/burnless/burnless/internal/config"
)

// NewValidateCmd creates the "burnless validate" command.
func NewValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "validate [file]",
		Short:        "Validate an sre.yaml file",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE:         runValidate,
	}
}

func runValidate(cmd *cobra.Command, args []string) error {
	path := "sre.yaml"
	if len(args) > 0 {
		path = args[0]
	}

	cfg, err := config.Load(path)
	if err != nil {
		fmt.Printf("✗ %s is invalid\n", path)
		return err
	}

	fmt.Printf("✓ %s is valid\n", path)
	fmt.Printf("  service: %s\n", cfg.Service)
	fmt.Printf("  %d SLO(s) defined\n", len(cfg.SLOs))
	fmt.Printf("  %d burn rate alert(s) configured\n", len(cfg.ErrorBudget.BurnRateAlerts))
	fmt.Printf("  %d runbook(s) defined\n", len(cfg.Runbooks))

	return nil
}
