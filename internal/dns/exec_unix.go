//go:build darwin || linux

package dns

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// sudo runs a command with elevated privileges, sharing the terminal so the
// password prompt is visible.
func sudo(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stderr, os.Stderr
	return cmd.Run()
}

// writeSudo writes content to a root-owned path via `sudo tee`.
func writeSudo(path, content string) error {
	if err := sudo("mkdir", "-p", filepath.Dir(path)); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "tee", path)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
