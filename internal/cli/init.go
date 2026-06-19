package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/burnless/burnless/internal/config"
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

	// validate before writing — fail fast if something is wrong
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("generated config is invalid: %w", err)
	}

	if err := config.WriteToFile(cfg, "sre.yaml"); err != nil {
		return fmt.Errorf("writing sre.yaml: %w", err)
	}

	fmt.Println("✓ Created sre.yaml")
	return nil
}

// prompt prints a question and reads a single line of user input.
func prompt(reader *bufio.Reader, question string) string {
	fmt.Printf("? %s: ", question)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}
