package cmd

import (
	"os"
	"path/filepath"

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

func mugenShellDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "quickshell", "mugen-shell")
}

func mugenShellPath(parts ...string) string {
	return filepath.Join(append([]string{mugenShellDir()}, parts...)...)
}
