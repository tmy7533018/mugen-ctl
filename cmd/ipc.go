package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var ipcCmd = &cobra.Command{
	Use:   "ipc <mode>",
	Short: "Send a command to mugen-shell via IPC",
	Long: `Send a command to the mugen-shell IPC file.

Available modes:
  launcher             Open app launcher
  calendar             Open calendar
  wallpaper            Open wallpaper selector
  music                Open music player
  notification         Open notifications
  powermenu            Open power menu
  volume               Open volume control
  window-switcher      Open window switcher
  window-switcher-next Select next window in switcher
  window-switcher-prev Select previous window in switcher
  close                Close all modules`,
	Args: cobra.MinimumNArgs(1),
	RunE: runIPC,
}

func init() {
	rootCmd.AddCommand(ipcCmd)
}

func ipcFilePath() string {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		dir = "/tmp"
	}
	return filepath.Join(dir, "mugen-shell-ipc")
}

func runIPC(_ *cobra.Command, args []string) error {
	msg := strings.Join(args, " ") + "\n"
	path := ipcFilePath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot create IPC dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open IPC file: %w", err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return fmt.Errorf("IPC busy (mugen-shell may be processing a command)")
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN) //nolint:errcheck

	_, err = fmt.Fprint(f, msg)
	return err
}
