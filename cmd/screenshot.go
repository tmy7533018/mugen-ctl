package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Take a region screenshot with slurp + grim",
	RunE:  runScreenshot,
}

func init() {
	rootCmd.AddCommand(screenshotCmd)
}

func runScreenshot(_ *cobra.Command, _ []string) error {
	script := mugenShellPath("scripts", "take-screenshot.sh")
	if _, err := os.Stat(script); err != nil {
		return fmt.Errorf("take-screenshot.sh not found at %s", script)
	}

	c := exec.Command("bash", script)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
