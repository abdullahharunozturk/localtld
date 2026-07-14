//go:build darwin

package service

import (
	"fmt"
	"strings"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

type darwinManager struct{}

func newManager() Manager { return darwinManager{} }

func label(name string) string     { return "sh.localtld." + name }
func plistPath(name string) string { return "/Library/LaunchDaemons/" + label(name) + ".plist" }

func (darwinManager) Install(u Unit) error {
	var args strings.Builder
	args.WriteString("<string>" + u.Exec + "</string>")
	for _, a := range u.Args {
		args.WriteString("\n    <string>" + a + "</string>")
	}
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>%s</string>
  <key>ProgramArguments</key><array>
    %s
  </array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
</dict></plist>
`, label(u.Name), args.String())
	if err := sysexec.WriteSudo(plistPath(u.Name), plist); err != nil {
		return err
	}
	// Reload cleanly whether or not it was already loaded.
	_ = sysexec.Sudo("launchctl", "bootout", "system", plistPath(u.Name))
	if err := sysexec.Sudo("launchctl", "bootstrap", "system", plistPath(u.Name)); err != nil {
		return err
	}
	return sysexec.Sudo("launchctl", "kickstart", "-k", "system/"+label(u.Name))
}

func (darwinManager) Uninstall(name string) error {
	_ = sysexec.Sudo("launchctl", "bootout", "system", plistPath(name))
	return sysexec.Sudo("rm", "-f", plistPath(name))
}
