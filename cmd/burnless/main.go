package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/burnless/burnless/internal/cli"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "burnless",
		Short: "Stop burning your error budget. Stop burning out your team.",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the burnless version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("burnless", version)
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cli.NewInitCmd())
	rootCmd.AddCommand(cli.NewValidateCmd())
	rootCmd.AddCommand(cli.NewToilCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
