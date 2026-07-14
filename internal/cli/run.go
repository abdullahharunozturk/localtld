package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/abdullahharunozturk/localtld/internal/config"
	"github.com/abdullahharunozturk/localtld/internal/netutil"
	"github.com/abdullahharunozturk/localtld/internal/proxy"
)

// Run executes a project's dev command under its localtld domain, falling back
// to a plain run whenever the domain can't be provided.
func Run(args []string) error {
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) == 0 {
		return errors.New("usage: localtld run -- <command>   (e.g. localtld run -- pnpm dev)")
	}

	cwd, _ := os.Getwd()
	label := config.ReadLabel(cwd)

	// Project doesn't opt in → run the command untouched.
	if label == "" {
		return execCommand(args, nil)
	}
	// Opted in but the machine isn't set up → localhost fallback.
	if !config.IsSetup() {
		fmt.Fprintln(os.Stderr, dim(`→ localtld not set up — running on localhost (run "localtld setup")`))
		return execCommand(args, nil)
	}

	host := label + "." + config.GetTLD()
	port, err := netutil.FreePort()
	if err != nil {
		return err
	}
	if err := proxy.AddSite(host, port); err != nil {
		return err
	}
	if err := proxy.Reload(); err != nil {
		_ = proxy.RemoveSite(host)
		fmt.Fprintln(os.Stderr, dim("→ caddy reload failed — running on localhost"))
		return execCommand(args, nil)
	}
	defer func() {
		_ = proxy.RemoveSite(host)
		_ = proxy.Reload()
	}()

	fmt.Fprintf(os.Stderr, "\n  %s  %s\n\n", green("→ http://"+host), dim(fmt.Sprintf("(port %d)", port)))
	return execCommand(args, []string{"PORT=" + fmt.Sprint(port)})
}

// execCommand runs the wrapped command, inheriting stdio and propagating its
// exit code.
func execCommand(args []string, extraEnv []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = append(os.Environ(), extraEnv...)
	err := cmd.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		os.Exit(ee.ExitCode())
	}
	return err
}
