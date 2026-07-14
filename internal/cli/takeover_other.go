//go:build !darwin

package cli

// freeConflictingServices is macOS-specific (brew launchd daemons); a no-op
// elsewhere.
func freeConflictingServices() {}
