package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var sddmCmd = &cobra.Command{
	Use:   "sddm",
	Short: "Manage SDDM theme assets",
}

var sddmRandomizeCmd = &cobra.Command{
	Use:   "randomize",
	Short: "Copy a random video from wallpapers/videos to the SDDM background",
	RunE:  runSddmRandomize,
}

func init() {
	rootCmd.AddCommand(sddmCmd)
	sddmCmd.AddCommand(sddmRandomizeCmd)
}

const sddmDestDir = "/usr/share/sddm/themes/sddm-astronaut-theme/Backgrounds"
const sddmTargetName = "login.mp4"

func runSddmRandomize(_ *cobra.Command, _ []string) error {
	if os.Geteuid() != 0 {
		c := exec.Command("sudo", append([]string{os.Args[0]}, os.Args[1:]...)...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	}

	srcDir, err := sddmVideoSourceDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("cannot read video dir %s: %w", srcDir, err)
	}

	var mp4s []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".mp4") {
			mp4s = append(mp4s, filepath.Join(srcDir, e.Name()))
		}
	}
	if len(mp4s) == 0 {
		return fmt.Errorf("no mp4 files found in %s", srcDir)
	}

	pick := mp4s[rand.Intn(len(mp4s))]

	if err := os.MkdirAll(sddmDestDir, 0755); err != nil {
		return err
	}

	destPath := filepath.Join(sddmDestDir, sddmTargetName)
	os.Remove(destPath) //nolint:errcheck

	if err := copyFile(pick, destPath, 0644); err != nil {
		return err
	}

	fmt.Printf("SDDM background set: %s\n", pick)
	return nil
}

func sddmVideoSourceDir() (string, error) {
	// When running under sudo, resolve the original user's home directory.
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		out, err := exec.Command("getent", "passwd", sudoUser).Output()
		if err == nil {
			parts := strings.Split(strings.TrimSpace(string(out)), ":")
			if len(parts) >= 6 {
				return filepath.Join(parts[5], ".config", "mugen-shell", "wallpapers", "videos"), nil
			}
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "quickshell", "mugen-shell", "wallpapers", "videos"), nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
