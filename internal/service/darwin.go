//go:build darwin

package service

import (
	"fmt"
	"os/exec"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

const darwinLabel = "sh.localtld.caddy"

type darwinManager struct{}

func newManager() Manager { return darwinManager{} }

func (darwinManager) EnsureCaddy(configPath string) error {
	caddy, err := exec.LookPath("caddy")
	if err != nil {
		return fmt.Errorf("caddy not found on PATH — install it (brew install caddy)")
	}
	plistPath := "/Library/LaunchDaemons/" + darwinLabel + ".plist"
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>%s</string>
  <key>ProgramArguments</key><array>
    <string>%s</string><string>run</string>
    <string>--config</string><string>%s</string>
    <string>--adapter</string><string>caddyfile</string>
  </array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
</dict></plist>
`, darwinLabel, caddy, configPath)
	if err := sysexec.WriteSudo(plistPath, plist); err != nil {
		return err
	}
	// Reload cleanly whether or not it was already loaded.
	_ = sysexec.Sudo("launchctl", "bootout", "system", plistPath)
	if err := sysexec.Sudo("launchctl", "bootstrap", "system", plistPath); err != nil {
		return err
	}
	return sysexec.Sudo("launchctl", "kickstart", "-k", "system/"+darwinLabel)
}

func (darwinManager) StopCaddy() error {
	plistPath := "/Library/LaunchDaemons/" + darwinLabel + ".plist"
	_ = sysexec.Sudo("launchctl", "bootout", "system", plistPath)
	return sysexec.Sudo("rm", "-f", plistPath)
}
