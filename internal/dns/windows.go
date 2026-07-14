//go:build windows

package dns

import (
	"fmt"
	"os"
	"os/exec"
)

type windowsProvider struct{}

func newProvider() Provider { return windowsProvider{} }

func (windowsProvider) Name() string { return "Windows (NRPT)" }

// ps runs a PowerShell command. NRPT changes require an elevated shell.
func ps(script string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	return cmd.Run()
}

func (p windowsProvider) Setup(tld string) error {
	// NRPT is the Windows equivalent of /etc/resolver: it routes a namespace to
	// a chosen resolver. A local DNS server on 127.0.0.1:53 (dnsmasq/CoreDNS)
	// must answer *.tld → 127.0.0.1 for this to resolve.
	ns := "." + tld
	if err := ps(fmt.Sprintf("Add-DnsClientNrptRule -Namespace '%s' -NameServers '127.0.0.1'", ns)); err != nil {
		return err
	}
	return p.FlushCache()
}

func (p windowsProvider) Teardown(tld string) error {
	ns := "." + tld
	_ = ps(fmt.Sprintf("Get-DnsClientNrptRule | Where-Object Namespace -eq '%s' | Remove-DnsClientNrptRule -Force", ns))
	return p.FlushCache()
}

func (windowsProvider) FlushCache() error {
	return exec.Command("ipconfig", "/flushdns").Run()
}
