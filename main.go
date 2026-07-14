// localtld — give your local projects real domains on dynamic ports.
//
// This Go binary is the cross-platform successor to bin/localtld (bash, macOS-only).
// The macOS path mirrors the bash tool; Linux and Windows add native DNS wiring
// (systemd-resolved / NRPT). See internal/dns for the per-OS providers.
package main

import (
	"fmt"
	"os"

	"github.com/abdullahharunozturk/localtld/internal/cli"
)

const version = "0.1.0"

func main() {
	// Never run under sudo: it would corrupt config ownership and package paths.
	// Privileged steps elevate themselves only when needed.
	if os.Geteuid() == 0 && os.Getenv("LOCALTLD_ALLOW_ROOT") == "" {
		fmt.Fprintln(os.Stderr, "localtld: do not run with sudo — it asks for admin rights only when it needs them")
		os.Exit(1)
	}
	if err := dispatch(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "localtld: "+err.Error())
		os.Exit(1)
	}
}

func dispatch(args []string) error {
	cmd := ""
	if len(args) > 0 {
		cmd = args[0]
	}
	switch cmd {
	case "", "help", "-h", "--help":
		cli.Usage(version)
		return nil
	case "version", "-v", "--version":
		fmt.Println("localtld " + version)
		return nil
	case "setup":
		return cli.Setup()
	case "run":
		return cli.Run(args[1:])
	case "list":
		return cli.List()
	case "tld":
		return cli.TLD(args[1:])
	case "doctor":
		return cli.Doctor()
	case "uninstall":
		return cli.Uninstall()
	default:
		return fmt.Errorf("unknown command %q (try: localtld help)", cmd)
	}
}
