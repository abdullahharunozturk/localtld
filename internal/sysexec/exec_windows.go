//go:build windows

// Package sysexec centralizes privileged/system command execution so the dns
// and service layers share one implementation per OS.
package sysexec

import (
	"os"
	"os/exec"
)

// PS runs a PowerShell command. Privileged operations assume an elevated shell.
func PS(script string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	return cmd.Run()
}
