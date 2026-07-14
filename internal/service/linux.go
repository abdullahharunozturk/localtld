//go:build linux

package service

import (
	"fmt"
	"os/exec"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

const linuxUnit = "localtld-caddy.service"

type linuxManager struct{}

func newManager() Manager { return linuxManager{} }

func (linuxManager) EnsureCaddy(configPath string) error {
	caddy, err := exec.LookPath("caddy")
	if err != nil {
		return fmt.Errorf("caddy not found on PATH — install it first")
	}
	// AmbientCapabilities lets Caddy bind :80 without running as full root.
	unit := fmt.Sprintf(`[Unit]
Description=localtld Caddy
After=network.target

[Service]
ExecStart=%s run --config %s --adapter caddyfile
Restart=on-failure
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
`, caddy, configPath)
	if err := sysexec.WriteSudo("/etc/systemd/system/"+linuxUnit, unit); err != nil {
		return err
	}
	if err := sysexec.Sudo("systemctl", "daemon-reload"); err != nil {
		return err
	}
	return sysexec.Sudo("systemctl", "enable", "--now", linuxUnit)
}

func (linuxManager) StopCaddy() error {
	_ = sysexec.Sudo("systemctl", "disable", "--now", linuxUnit)
	_ = sysexec.Sudo("rm", "-f", "/etc/systemd/system/"+linuxUnit)
	return sysexec.Sudo("systemctl", "daemon-reload")
}
