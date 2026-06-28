package cli

import (
	"fmt"
	"time"

	"github.com/burnless/burnless/internal/toil"

	"github.com/spf13/cobra"
)

// NewToilCmd creates the "burnless toil" parent command.
func NewToilCmd() *cobra.Command {
	toilCmd := &cobra.Command{
		Use:   "toil",
		Short: "Track and report on manual SRE toil",
	}

	toilCmd.AddCommand(newToilLogCmd())
	toilCmd.AddCommand(newToilReportCmd())
	toilCmd.AddCommand(newToilExportCmd())

	return toilCmd
}

func newToilLogCmd() *cobra.Command {
	var (
		service     string
		task        string
		duration    string
		automatable bool
		trigger     string
		notes       string
	)

	cmd := &cobra.Command{
		Use:          "log",
		Short:        "Log a toil event",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runToilLog(service, task, duration, automatable, trigger, notes)
		},
	}

	cmd.Flags().StringVar(&service, "service", "", "Service name (required)")
	cmd.Flags().StringVar(&task, "task", "", "Task name, e.g. manual-rollback (required)")
	cmd.Flags().StringVar(&duration, "duration", "", "Time spent, e.g. 45m (required)")
	cmd.Flags().BoolVar(&automatable, "automatable", false, "Could this task be automated?")
	cmd.Flags().StringVar(&trigger, "trigger", "", "What triggered this task")
	cmd.Flags().StringVar(&notes, "notes", "", "Free-text notes")

	if err := cmd.MarkFlagRequired("service"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("task"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("duration"); err != nil {
		panic(err)
	}

	return cmd
}

func runToilLog(service, task, duration string, automatable bool, trigger, notes string) error {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("invalid --duration %q: %w (try something like 45m or 1h30m)", duration, err)
	}

	event := toil.Event{
		ID:           fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Service:      service,
		Task:         task,
		Date:         time.Now(),
		DurationMins: int(d.Minutes()),
		Trigger:      trigger,
		Automatable:  automatable,
		Notes:        notes,
	}

	dbPath, err := toil.DefaultDBPath()
	if err != nil {
		return fmt.Errorf("finding database path: %w", err)
	}

	store, err := toil.OpenStore(dbPath)
	if err != nil {
		return fmt.Errorf("opening toil database: %w", err)
	}
	defer func() { _ = store.Close() }()

	if err := store.Save(event); err != nil {
		return fmt.Errorf("saving toil event: %w", err)
	}

	fmt.Printf("✓ Logged %d min of toil for %s (%s)\n", event.DurationMins, service, task)
	return nil
}
