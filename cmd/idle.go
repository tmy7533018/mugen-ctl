package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var idleCmd = &cobra.Command{
	Use:   "idle",
	Short: "Manage the idle inhibitor",
}

var idleToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle hypridle.service on/off",
	RunE:  runIdleToggle,
}

func init() {
	rootCmd.AddCommand(idleCmd)
	idleCmd.AddCommand(idleToggleCmd)
}

func runIdleToggle(_ *cobra.Command, _ []string) error {
	isActive := exec.Command("systemctl", "--user", "is-active", "--quiet", "hypridle.service").Run() == nil

	var action string
	if isActive {
		action = "stop"
	} else {
		action = "start"
	}

	out, err := exec.Command("systemctl", "--user", action, "hypridle.service").CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %s hypridle: %s", action, string(out))
	}

	if isActive {
		fmt.Println("hypridle stopped (idle inhibited)")
	} else {
		fmt.Println("hypridle started (idle active)")
	}
	return nil
}
