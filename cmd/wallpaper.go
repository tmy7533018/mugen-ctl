package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var wallpaperCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Manage the desktop wallpaper",
}

var wallpaperSetCmd = &cobra.Command{
	Use:   "set <path>",
	Short: "Set a wallpaper from an image or video file",
	Args:  cobra.ExactArgs(1),
	RunE:  runWallpaperSet,
}

var wallpaperGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Print the current wallpaper path",
	RunE:  runWallpaperGet,
}

func init() {
	rootCmd.AddCommand(wallpaperCmd)
	wallpaperCmd.AddCommand(wallpaperSetCmd, wallpaperGetCmd)
}

func runWallpaperSet(_ *cobra.Command, args []string) error {
	script := mugenShellPath("scripts", "change-wallpaper.sh")
	if _, err := os.Stat(script); err != nil {
		return fmt.Errorf("change-wallpaper.sh not found at %s", script)
	}

	c := exec.Command("bash", script, args[0])
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func runWallpaperGet(_ *cobra.Command, _ []string) error {
	p := mugenShellPath(".cache", "wallp", "current_wallpaper_path.txt")
	data, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("no wallpaper recorded (cache not found)")
	}
	fmt.Print(string(data))
	return nil
}
