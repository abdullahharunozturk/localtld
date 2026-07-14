// Package proxy manages the Caddy config that maps hosts to dynamic ports.
package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/abdullahharunozturk/localtld/internal/config"
)

// WriteBase writes the root Caddyfile that imports every per-site snippet.
func WriteBase() error {
	if err := os.MkdirAll(config.SitesDir(), 0o755); err != nil {
		return err
	}
	base := fmt.Sprintf("{\n\tauto_https off\n}\n\nimport %s\n",
		filepath.Join(config.SitesDir(), "*.caddy"))
	return os.WriteFile(config.CaddyfilePath(), []byte(base), 0o644)
}

// SitePath is the snippet file for a host.
func SitePath(host string) string {
	return filepath.Join(config.SitesDir(), host+".caddy")
}

// AddSite maps http://host → 127.0.0.1:port.
func AddSite(host string, port int) error {
	if err := os.MkdirAll(config.SitesDir(), 0o755); err != nil {
		return err
	}
	snip := fmt.Sprintf("http://%s {\n\treverse_proxy 127.0.0.1:%d\n}\n", host, port)
	return os.WriteFile(SitePath(host), []byte(snip), 0o644)
}

// RemoveSite drops a host mapping (no error if already gone).
func RemoveSite(host string) error {
	if err := os.Remove(SitePath(host)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Reload tells the running Caddy to pick up the current config.
func Reload() error {
	cmd := exec.Command("caddy", "reload",
		"--config", config.CaddyfilePath(), "--adapter", "caddyfile")
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
