//go:build darwin || linux

// Package sysexec centralizes privileged/system command execution so the dns
// and service layers share one implementation per OS.
package sysexec

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Sudo runs a command with elevated privileges, sharing the terminal so the
// password prompt is visible.
func Sudo(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stderr, os.Stderr
	return cmd.Run()
}

// WriteSudo writes content to a root-owned path via `sudo tee`.
func WriteSudo(path, content string) error {
	if err := Sudo("mkdir", "-p", filepath.Dir(path)); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "tee", path)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
