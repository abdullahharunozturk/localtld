// Package cli implements the localtld subcommands.
package cli

import (
	"fmt"
	"os"

	"github.com/abdullahharunozturk/localtld/internal/config"
)

var useColor = os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"

func paint(code, s string) string {
	if !useColor {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func dim(s string) string   { return paint("2", s) }
func bold(s string) string  { return paint("1", s) }
func green(s string) string { return paint("32", s) }

// Usage prints the help text.
func Usage(version string) {
	fmt.Printf(`localtld %s — give your local projects a domain

Commands
  setup                First-time setup: choose a TLD + configure DNS/Caddy
  run -- <cmd>         Run a project under its domain (falls back if not set up)
  list                 Show active projects and their domains
  tld <new>            Change the machine-wide TLD
  doctor               Check the health of your setup
  uninstall            Revert DNS/route changes
  version              Version

Current TLD  .%s   (change: localtld tld <new>)
`, version, config.GetTLD())
}
