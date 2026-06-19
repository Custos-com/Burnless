package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/burnless/burnless/internal/config"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the "burnless init" command.
func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a new sre.yaml for your service",
		RunE:  runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	serviceName := prompt(reader, "Service name")
	team := prompt(reader, "Team")

	cfg := config.GenerateDefault(serviceName, team)

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("generated config is invalid: %w", err)
	}

	if err := config.WriteToFile(cfg, "sre.yaml"); err != nil {
		return fmt.Errorf("writing sre.yaml: %w", err)
	}

	fmt.Println("✓ Created sre.yaml")
	return nil
}

func prompt(reader *bufio.Reader, question string) string {
	fmt.Printf("? %s: ", question)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}
