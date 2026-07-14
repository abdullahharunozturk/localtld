//go:build darwin

package dns

import (
	"path/filepath"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

type darwinProvider struct{}

func newProvider() Provider { return darwinProvider{} }

func (darwinProvider) Name() string { return "macOS (/etc/resolver)" }

func (p darwinProvider) Setup(tld string) error {
	// /etc/resolver/<tld> makes macOS send *.tld queries to 127.0.0.1:53,
	// where localtld's built-in DNS server answers.
	if err := sysexec.WriteSudo(filepath.Join("/etc/resolver", tld), "nameserver 127.0.0.1\n"); err != nil {
		return err
	}
	return p.FlushCache()
}

func (p darwinProvider) Teardown(tld string) error {
	_ = sysexec.Sudo("rm", "-f", filepath.Join("/etc/resolver", tld))
	return p.FlushCache()
}

func (darwinProvider) FlushCache() error {
	_ = sysexec.Sudo("dscacheutil", "-flushcache")
	// macOS 26 (Tahoe) needs a full restart of mDNSResponder, not just -HUP.
	return sysexec.Sudo("killall", "mDNSResponder")
}
