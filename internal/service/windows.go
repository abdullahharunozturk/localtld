//go:build windows

package service

import (
	"os"
	"os/exec"
)

type windowsManager struct{}

func newManager() Manager { return windowsManager{} }

// On Windows, binding :80 doesn't require elevation, so Caddy runs in the
// background via its own `start`/`stop` lifecycle.
func (windowsManager) EnsureCaddy(configPath string) error {
	cmd := exec.Command("caddy", "start", "--config", configPath, "--adapter", "caddyfile")
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	return cmd.Run()
}

func (windowsManager) StopCaddy() error {
	cmd := exec.Command("caddy", "stop")
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	return cmd.Run()
}
