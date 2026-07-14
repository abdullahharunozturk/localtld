// Package dns wires "*.<tld> → 127.0.0.1" into the host OS resolver.
//
// Each OS has its own mechanism, selected at build time:
//
//	macOS    /etc/resolver/<tld>          (darwin.go)
//	Linux    systemd-resolved + dnsmasq   (linux.go)
//	Windows  NRPT rule                    (windows.go)
package dns

// Provider is the per-OS resolver integration.
type Provider interface {
	// Name is a short human label for `doctor`/`setup` output.
	Name() string
	// Setup routes *.tld queries to the local resolver.
	Setup(tld string) error
	// Teardown removes the routing for tld.
	Teardown(tld string) error
	// FlushCache clears the OS DNS cache so changes take effect immediately.
	FlushCache() error
}

// New returns the provider for the current OS (see the build-tagged files).
func New() Provider { return newProvider() }
