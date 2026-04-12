package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var blurCmd = &cobra.Command{
	Use:   "blur",
	Short: "Manage Hyprland blur presets",
}

var blurListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available blur presets",
	RunE:  runBlurList,
}

var blurCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print the active blur preset",
	RunE:  runBlurCurrent,
}

var blurSetCmd = &cobra.Command{
	Use:   "set <preset>",
	Short: "Apply a blur preset by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runBlurSet,
}

func init() {
	rootCmd.AddCommand(blurCmd)
	blurCmd.AddCommand(blurListCmd, blurCurrentCmd, blurSetCmd)
}

func looknfeelConf() string {
	if v := os.Getenv("HYPR_CONFIG_FILE"); v != "" {
		return v
	}
	cfg := os.Getenv("XDG_CONFIG_HOME")
	if cfg == "" {
		home, _ := os.UserHomeDir()
		cfg = filepath.Join(home, ".config")
	}
	return filepath.Join(cfg, "hypr", "configs", "looknfeel.conf")
}

var (
	blurBlockRe    = regexp.MustCompile(`(^|[\s])blur\s*\{`)
	presetHeaderRe = regexp.MustCompile(`^[\t ]*#[\t ]*={3,}[\t ]+(.+?)[\t ]+={3,}[\t ]*$`)
	assignRe       = regexp.MustCompile(`[a-z_]+\s*=`)
)

type blurPreset struct {
	name  string
	lines []int // indices of assignment lines within allLines
}

func parseBlurPresets(conf string) ([]blurPreset, []string, error) {
	f, err := os.Open(conf)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var allLines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		allLines = append(allLines, sc.Text())
	}

	inBlur := false
	var presets []blurPreset
	var current *blurPreset

	for i, line := range allLines {
		if !inBlur {
			if blurBlockRe.MatchString(line) {
				inBlur = true
			}
			continue
		}
		if strings.Contains(line, "}") {
			if current != nil {
				presets = append(presets, *current)
				current = nil
			}
			inBlur = false
			continue
		}
		if m := presetHeaderRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				presets = append(presets, *current)
			}
			current = &blurPreset{name: strings.TrimSpace(m[1])}
			continue
		}
		if current != nil && assignRe.MatchString(line) {
			current.lines = append(current.lines, i)
		}
	}

	return presets, allLines, nil
}

func runBlurList(_ *cobra.Command, _ []string) error {
	presets, _, err := parseBlurPresets(looknfeelConf())
	if err != nil {
		return err
	}
	for _, p := range presets {
		fmt.Println(p.name)
	}
	return nil
}

func runBlurCurrent(_ *cobra.Command, _ []string) error {
	presets, allLines, err := parseBlurPresets(looknfeelConf())
	if err != nil {
		return err
	}
	for _, p := range presets {
		for _, li := range p.lines {
			if !strings.HasPrefix(strings.TrimSpace(allLines[li]), "#") {
				fmt.Println(p.name)
				return nil
			}
		}
	}
	fmt.Println("(none)")
	return nil
}

func runBlurSet(_ *cobra.Command, args []string) error {
	want := args[0]
	conf := looknfeelConf()

	presets, allLines, err := parseBlurPresets(conf)
	if err != nil {
		return err
	}

	target := ""
	for _, p := range presets {
		if strings.EqualFold(p.name, want) || strings.Contains(strings.ToLower(p.name), strings.ToLower(want)) {
			target = p.name
			break
		}
	}
	if target == "" {
		names := make([]string, len(presets))
		for i, p := range presets {
			names[i] = p.name
		}
		return fmt.Errorf("preset %q not found\navailable: %s", want, strings.Join(names, ", "))
	}

	modified := make([]string, len(allLines))
	copy(modified, allLines)

	for _, p := range presets {
		isTarget := p.name == target
		for _, li := range p.lines {
			line := allLines[li]
			trimmed := strings.TrimSpace(line)
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]

			if isTarget {
				if strings.HasPrefix(trimmed, "#") {
					rest := strings.TrimLeft(strings.TrimPrefix(trimmed, "#"), " ")
					modified[li] = indent + rest
				}
			} else {
				if !strings.HasPrefix(trimmed, "#") {
					modified[li] = indent + "# " + trimmed
				}
			}
		}
	}

	tmp := conf + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(out)
	for _, l := range modified {
		fmt.Fprintln(w, l)
	}
	if err := w.Flush(); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	out.Close()

	if err := os.Rename(tmp, conf); err != nil {
		os.Remove(tmp)
		return err
	}

	fmt.Printf("Active preset: %s\n", target)

	if _, err := exec.LookPath("hyprctl"); err == nil {
		exec.Command("hyprctl", "reload").Run() //nolint:errcheck
		fmt.Println("Hyprland config reloaded")
	}

	return nil
}
