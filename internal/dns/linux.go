//go:build linux

package dns

import (
	"fmt"
	"path/filepath"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

type linuxProvider struct{}

func newProvider() Provider { return linuxProvider{} }

func (linuxProvider) Name() string { return "Linux (systemd-resolved)" }

func (p linuxProvider) Setup(tld string) error {
	// Route *.tld to localtld's built-in DNS server (127.0.0.1:53) via a
	// systemd-resolved drop-in.
	drop := filepath.Join("/etc/systemd/resolved.conf.d", "localtld.conf")
	if err := sysexec.WriteSudo(drop, fmt.Sprintf("[Resolve]\nDNS=127.0.0.1\nDomains=~%s\n", tld)); err != nil {
		return err
	}
	if err := sysexec.Sudo("systemctl", "restart", "systemd-resolved"); err != nil {
		return err
	}
	return p.FlushCache()
}

func (p linuxProvider) Teardown(tld string) error {
	_ = sysexec.Sudo("rm", "-f", "/etc/systemd/resolved.conf.d/localtld.conf")
	_ = sysexec.Sudo("systemctl", "restart", "systemd-resolved")
	return p.FlushCache()
}

func (linuxProvider) FlushCache() error {
	return sysexec.Sudo("resolvectl", "flush-caches")
}
