package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// promptTLD lets the user pick a TLD interactively, defaulting to current.
func promptTLD(current string) string {
	if !isTTY(os.Stdin) {
		return current
	}
	fmt.Fprintf(os.Stderr, "  TLD to use [%s]: ", current)
	var ans string
	_, _ = fmt.Fscanln(os.Stdin, &ans)
	if ans = strings.TrimPrefix(strings.TrimSpace(ans), "."); ans != "" {
		return ans
	}
	return current
}

// confirmYesNo asks a [y/N] question; a non-TTY declines.
func confirmYesNo(q string) bool {
	if !isTTY(os.Stdin) {
		return false
	}
	fmt.Fprintf(os.Stderr, "  %s [y/N] ", q)
	var ans string
	_, _ = fmt.Fscanln(os.Stdin, &ans)
	a := strings.ToLower(strings.TrimSpace(ans))
	return a == "y" || a == "yes"
}

func runVisible(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stderr, os.Stderr
	return cmd.Run()
}

// ensureCaddy returns the caddy path, offering to install it when missing.
func ensureCaddy() (string, error) {
	if p, err := exec.LookPath("caddy"); err == nil {
		return p, nil
	}
	switch runtime.GOOS {
	case "darwin":
		if confirmYesNo("caddy is not installed. Install it with Homebrew now?") {
			if err := runVisible("brew", "install", "caddy"); err != nil {
				return "", err
			}
			return exec.LookPath("caddy")
		}
		return "", errors.New("caddy is required — install it: brew install caddy")
	case "windows":
		return "", errors.New("caddy is required — install it: winget install CaddyServer.Caddy")
	default:
		return "", errors.New("caddy is required — see https://caddyserver.com/docs/install")
	}
}
