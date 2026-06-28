package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/burnless/burnless/internal/toil"

	"github.com/spf13/cobra"
)

// NewToilExportCmd creates the "burnless toil export" command.
func newToilExportCmd() *cobra.Command {
	var (
		format string
		month  string
		output string
	)

	cmd := &cobra.Command{
		Use:          "export",
		Short:        "Export toil events as CSV or JSON",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runToilExport(format, month, output)
		},
	}

	cmd.Flags().StringVar(&format, "format", "csv", "Output format: csv or json")
	cmd.Flags().StringVar(&month, "month", "", "Filter by month e.g. 2026-06 (default: all)")
	cmd.Flags().StringVar(&output, "output", "", "File path to write to (default: stdout)")

	return cmd
}

func runToilExport(format, month, output string) error {
	dbPath, err := toil.DefaultDBPath()
	if err != nil {
		return fmt.Errorf("finding database path: %w", err)
	}

	store, err := toil.OpenStore(dbPath)
	if err != nil {
		return fmt.Errorf("opening toil database: %w", err)
	}
	defer func() { _ = store.Close() }()

	events, err := store.All()
	if err != nil {
		return fmt.Errorf("reading toil events: %w", err)
	}

	// filter by month if specified
	if month != "" {
		events, err = toil.FilterByMonth(events, month)
		if err != nil {
			return err
		}
	}

	if len(events) == 0 {
		fmt.Fprintln(os.Stderr, "No toil events found.")
		return nil
	}

	// determine output writer
	var w io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer func() { _ = f.Close() }()
		w = f
	}

	model := toil.DefaultCostModel()

	switch format {
	case "csv":
		return exportCSV(w, events, model)
	case "json":
		return exportJSON(w, events, model)
	default:
		return fmt.Errorf("unknown format %q — use csv or json", format)
	}
}

func exportCSV(w io.Writer, events []toil.Event, model toil.CostModel) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	// header
	if err := cw.Write([]string{
		"id", "service", "task", "date",
		"duration_minutes", "automatable", "trigger", "notes", "cost_usd",
	}); err != nil {
		return err
	}

	// rows
	for _, e := range events {
		cost := model.Cost(e)
		if err := cw.Write([]string{
			e.ID,
			e.Service,
			e.Task,
			e.Date.Format("2006-01-02"),
			strconv.Itoa(e.DurationMins),
			strconv.FormatBool(e.Automatable),
			e.Trigger,
			e.Notes,
			fmt.Sprintf("%.2f", cost),
		}); err != nil {
			return err
		}
	}

	return cw.Error()
}

type exportRow struct {
	toil.Event
	CostUSD float64 `json:"cost_usd"`
}

func exportJSON(w io.Writer, events []toil.Event, model toil.CostModel) error {
	rows := make([]exportRow, len(events))
	for i, e := range events {
		rows[i] = exportRow{
			Event:   e,
			CostUSD: model.Cost(e),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
