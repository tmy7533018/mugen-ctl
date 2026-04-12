package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
	dir := os.Getenv("XDG_PICTURES_DIR")
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, "Pictures")
	}
	outDir := filepath.Join(dir, "mugen-screenshots")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	outPath := filepath.Join(outDir, "screenshot_"+time.Now().Format("20060102_150405")+".png")

	regionBytes, err := exec.Command("slurp").Output()
	if err != nil {
		return fmt.Errorf("slurp cancelled or failed")
	}
	region := strings.TrimSpace(string(regionBytes))

	c := exec.Command("grim", "-g", region, outPath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("grim failed: %w", err)
	}

	fmt.Println(outPath)
	return nil
}
