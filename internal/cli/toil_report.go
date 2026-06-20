package cli

import (
	"fmt"

	"github.com/burnless/burnless/internal/toil"

	"github.com/spf13/cobra"
)

func newToilReportCmd() *cobra.Command {
	var month string

	cmd := &cobra.Command{
		Use:          "report",
		Short:        "Generate a monthly toil report",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runToilReport(month)
		},
	}

	cmd.Flags().StringVar(&month, "month", "", "Month to report on, e.g. 2026-06 (required)")
	if err := cmd.MarkFlagRequired("month"); err != nil {
		panic(err)
	}

	return cmd
}

func runToilReport(month string) error {
	dbPath, err := toil.DefaultDBPath()
	if err != nil {
		return fmt.Errorf("finding database path: %w", err)
	}

	store, err := toil.OpenStore(dbPath)
	if err != nil {
		return fmt.Errorf("opening toil database: %w", err)
	}
	defer func() { _ = store.Close() }()

	allEvents, err := store.All()
	if err != nil {
		return fmt.Errorf("reading toil events: %w", err)
	}

	events, err := toil.FilterByMonth(allEvents, month)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		fmt.Printf("No toil events logged for %s\n", month)
		return nil
	}

	totalMinutes := toil.TotalMinutes(events)
	totalHours := float64(totalMinutes) / 60.0

	model := toil.DefaultCostModel()
	totalCost := model.CostForMinutes(totalMinutes)

	summaries := toil.GroupByTask(events, model)

	fmt.Printf("┌──────────────────────────────────────────────\n")
	fmt.Printf("│  Toil Report · %s\n", month)
	fmt.Printf("│\n")
	fmt.Printf("│  Total toil time     %.1f hours\n", totalHours)
	fmt.Printf("│  Estimated cost      $%.2f\n", totalCost)
	fmt.Printf("│  Events logged       %d\n", len(events))
	fmt.Printf("│\n")
	fmt.Printf("│  Top tasks by time:\n")
	for i, s := range summaries {
		if i >= 5 {
			break // show top 5 only
		}
		auto := ""
		if s.Automatable {
			auto = " (automatable)"
		}
		fmt.Printf("│  %d. %-20s %5.1f hrs  → %s%s\n", i+1, s.Task, s.TotalHours, s.Priority, auto)
	}
	fmt.Printf("└──────────────────────────────────────────────\n")

	return nil
}
