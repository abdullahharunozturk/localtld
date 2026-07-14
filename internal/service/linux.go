//go:build linux

package service

import (
	"fmt"
	"strings"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

type linuxManager struct{}

func newManager() Manager { return linuxManager{} }

func unitName(name string) string { return "localtld-" + name + ".service" }
func unitPath(name string) string { return "/etc/systemd/system/" + unitName(name) }

func (linuxManager) Install(u Unit) error {
	execStart := u.Exec
	if len(u.Args) > 0 {
		execStart += " " + strings.Join(u.Args, " ")
	}
	// AmbientCapabilities lets the service bind privileged ports (:53, :80)
	// without running as full root.
	unit := fmt.Sprintf(`[Unit]
Description=localtld %s
After=network.target

[Service]
ExecStart=%s
Restart=on-failure
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
`, u.Name, execStart)
	if err := sysexec.WriteSudo(unitPath(u.Name), unit); err != nil {
		return err
	}
	if err := sysexec.Sudo("systemctl", "daemon-reload"); err != nil {
		return err
	}
	return sysexec.Sudo("systemctl", "enable", "--now", unitName(u.Name))
}

func (linuxManager) Uninstall(name string) error {
	_ = sysexec.Sudo("systemctl", "disable", "--now", unitName(name))
	_ = sysexec.Sudo("rm", "-f", unitPath(name))
	return sysexec.Sudo("systemctl", "daemon-reload")
}
