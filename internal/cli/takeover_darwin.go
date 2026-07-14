//go:build darwin

package cli

import (
	"fmt"
	"os"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

// freeConflictingServices stops brew-managed dnsmasq/caddy that would fight
// localtld for :53/:80 — e.g. left over from the bash version. localtld needs
// those ports exclusively, so setup takes them over (a no-op on a fresh Mac).
func freeConflictingServices() {
	for _, svc := range []string{"dnsmasq", "caddy"} {
		plist := "/Library/LaunchDaemons/homebrew.mxcl." + svc + ".plist"
		if _, err := os.Stat(plist); err != nil {
			continue
		}
		fmt.Fprintf(os.Stderr, "  → stopping brew %s (it holds a port localtld needs)\n", svc)
		// bootout (not `brew services stop`) avoids running brew as root.
		_ = sysexec.SudoQuiet("launchctl", "bootout", "system", plist)
	}
}
