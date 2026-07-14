package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

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
	// Opted in but not set up → offer setup on an interactive terminal,
	// otherwise fall back to localhost.
	if !config.IsSetup() {
		if !offerSetup(label) {
			return execCommand(args, nil)
		}
		if err := Setup(); err != nil {
			fmt.Fprintln(os.Stderr, dim("→ setup did not complete — running on localhost"))
			return execCommand(args, nil)
		}
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
	cleanup := func() {
		_ = proxy.RemoveSite(host)
		_ = proxy.Reload()
	}
	// Ensure cleanup runs even on Ctrl-C: catch the signal so we survive long
	// enough to tidy the snippet. The child also receives it from the terminal
	// and shuts itself down, which returns cmd.Run below. (execCommand's os.Exit
	// would otherwise skip a deferred cleanup entirely.)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigc)

	fmt.Fprintf(os.Stderr, "\n  %s  %s\n\n", green("→ http://"+host), dim(fmt.Sprintf("(port %d)", port)))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	runErr := cmd.Run()

	cleanup()
	var ee *exec.ExitError
	if errors.As(runErr, &ee) {
		os.Exit(ee.ExitCode())
	}
	return runErr
}

// offerSetup asks the user (on an interactive terminal) whether to run first-time
// setup. On a non-TTY (CI, pipe) it silently declines and notes the localhost run.
func offerSetup(label string) bool {
	if !isTTY(os.Stdin) || !isTTY(os.Stderr) {
		fmt.Fprintln(os.Stderr, dim(`→ localtld not set up — running on localhost (run "localtld setup")`))
		return false
	}
	fmt.Fprintf(os.Stderr, "\n%s Setting up the domain (%s) needs your password once.\n",
		bold("→ Running localtld for the first time."), green(label+"."+config.GetTLD()))
	fmt.Fprint(os.Stderr, "  Set it up now? [y/N] ")
	var ans string
	_, _ = fmt.Fscanln(os.Stdin, &ans)
	switch strings.ToLower(strings.TrimSpace(ans)) {
	case "y", "yes":
		return true
	default:
		fmt.Fprintln(os.Stderr, dim(`  skipped — running on localhost (run "localtld setup" anytime)`))
		return false
	}
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
