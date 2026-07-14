// Package service installs long-running localtld background services (the
// built-in DNS server on :53 and Caddy on :80) as native OS services:
//
//	macOS    launchd LaunchDaemon (root)          darwin.go
//	Linux    systemd unit + CAP_NET_BIND_SERVICE  linux.go
//	Windows  scheduled task (highest privileges)  windows.go
package service

// Unit is a background service to install.
type Unit struct {
	Name string   // short id, e.g. "dns" or "caddy"
	Exec string   // absolute path to the executable
	Args []string // arguments
}

// Manager installs and removes OS services.
type Manager interface {
	Install(u Unit) error
	Uninstall(name string) error
}

// New returns the service manager for the current OS.
func New() Manager { return newManager() }
