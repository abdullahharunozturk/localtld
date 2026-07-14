// Package service manages the background Caddy process that binds :80.
//
// Binding a privileged port differs per OS:
//
//	macOS    launchd LaunchDaemon (root)          darwin.go
//	Linux    systemd unit + CAP_NET_BIND_SERVICE  linux.go
//	Windows  caddy start (no elevation needed)     windows.go
package service

// Manager controls the privileged Caddy process.
type Manager interface {
	// EnsureCaddy installs/starts a background Caddy bound to configPath.
	EnsureCaddy(configPath string) error
	// StopCaddy stops and removes the service.
	StopCaddy() error
}

// New returns the service manager for the current OS.
func New() Manager { return newManager() }
