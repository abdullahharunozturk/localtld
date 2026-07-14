//go:build darwin

package dns

import (
	"fmt"
	"os"
	"path/filepath"
)

type darwinProvider struct{}

func newProvider() Provider { return darwinProvider{} }

func (darwinProvider) Name() string { return "macOS (/etc/resolver + dnsmasq)" }

func (p darwinProvider) Setup(tld string) error {
	// dnsmasq answers *.tld → 127.0.0.1 (Homebrew prefix, Apple Silicon or Intel).
	prefix := "/opt/homebrew"
	if _, err := os.Stat(prefix); err != nil {
		prefix = "/usr/local"
	}
	conf := filepath.Join(prefix, "etc", "dnsmasq.d", "localtld.conf")
	if err := writeSudo(conf, fmt.Sprintf("address=/%s/127.0.0.1\n", tld)); err != nil {
		return err
	}
	// /etc/resolver/<tld> makes macOS send *.tld queries to the local resolver.
	if err := writeSudo(filepath.Join("/etc/resolver", tld), "nameserver 127.0.0.1\n"); err != nil {
		return err
	}
	_ = sudo("brew", "services", "restart", "dnsmasq")
	return p.FlushCache()
}

func (p darwinProvider) Teardown(tld string) error {
	_ = sudo("rm", "-f", filepath.Join("/etc/resolver", tld))
	return p.FlushCache()
}

func (darwinProvider) FlushCache() error {
	_ = sudo("dscacheutil", "-flushcache")
	// macOS 26 (Tahoe) needs a full restart of mDNSResponder, not just -HUP.
	return sudo("killall", "mDNSResponder")
}
