package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mugen-ctl",
	Short: "Control interface for mugen-shell",
	Long:  "mugen-ctl manages mugen-shell components running on Hyprland.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
