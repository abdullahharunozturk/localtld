//go:build windows

package service

import (
	"fmt"
	"strings"

	"github.com/abdullahharunozturk/localtld/internal/sysexec"
)

type windowsManager struct{}

func newManager() Manager { return windowsManager{} }

func taskName(name string) string { return "localtld-" + name }

// Install registers a scheduled task that runs the command at logon with the
// highest privileges (needed to bind :53/:80), then starts it now.
func (windowsManager) Install(u Unit) error {
	cmd := u.Exec
	if len(u.Args) > 0 {
		cmd += " " + strings.Join(u.Args, " ")
	}
	create := fmt.Sprintf(`schtasks /Create /TN "%s" /TR '%s' /SC ONLOGON /RL HIGHEST /F`, taskName(u.Name), cmd)
	if err := sysexec.PS(create); err != nil {
		return err
	}
	return sysexec.PS(fmt.Sprintf(`schtasks /Run /TN "%s"`, taskName(u.Name)))
}

func (windowsManager) Uninstall(name string) error {
	_ = sysexec.PS(fmt.Sprintf(`schtasks /End /TN "%s"`, taskName(name)))
	return sysexec.PS(fmt.Sprintf(`schtasks /Delete /TN "%s" /F`, taskName(name)))
}
