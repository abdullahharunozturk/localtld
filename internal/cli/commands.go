package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/abdullahharunozturk/localtld/internal/config"
	"github.com/abdullahharunozturk/localtld/internal/dns"
	"github.com/abdullahharunozturk/localtld/internal/proxy"
	"github.com/abdullahharunozturk/localtld/internal/resolver"
	"github.com/abdullahharunozturk/localtld/internal/service"
)

// Setup configures DNS + Caddy for the current machine TLD.
func Setup() error {
	tld := config.GetTLD()
	fmt.Println(bold("localtld setup"))
	fmt.Printf("  TLD: .%s\n", tld)
	if err := config.SetTLD(tld); err != nil {
		return err
	}
	if err := apply(tld); err != nil {
		return err
	}
	fmt.Println(green("✓ setup complete"))
	fmt.Println(dim(`  Add "localtld": "myapp" to package.json, then: localtld run -- <dev cmd>`))
	return nil
}

// apply installs the background services (built-in DNS + Caddy) and wires the
// OS resolver for tld. Used by both setup and tld-change.
func apply(tld string) error {
	if err := proxy.WriteBase(); err != nil {
		return err
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	caddy, err := exec.LookPath("caddy")
	if err != nil {
		return fmt.Errorf("caddy not found on PATH — install it (brew install caddy)")
	}
	svc := service.New()
	fmt.Println("  DNS server: installing service on :53 (needs admin once)")
	if err := svc.Install(service.Unit{Name: "dns", Exec: self, Args: []string{"serve-dns", tld}}); err != nil {
		return fmt.Errorf("dns service: %w", err)
	}
	fmt.Println("  Caddy: installing service on :80")
	if err := svc.Install(service.Unit{Name: "caddy", Exec: caddy,
		Args: []string{"run", "--config", config.CaddyfilePath(), "--adapter", "caddyfile"}}); err != nil {
		return fmt.Errorf("caddy service: %w", err)
	}
	p := dns.New()
	fmt.Printf("  Resolver: %s\n", p.Name())
	if err := p.Setup(tld); err != nil {
		return fmt.Errorf("resolver wiring: %w", err)
	}
	return nil
}

// ServeDNS runs the built-in DNS server (hidden command invoked by the service).
// The TLD is passed as an argument because the service runs as root and can't
// read the user's config.
func ServeDNS(args []string) error {
	tld := config.GetTLD()
	if len(args) > 0 && args[0] != "" {
		tld = args[0]
	}
	return resolver.Serve("127.0.0.1:53", tld)
}

// TLD changes the machine-wide TLD, re-wiring services if already set up.
func TLD(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: localtld tld <new-tld>   (e.g. localtld tld test)")
	}
	newTLD := strings.TrimPrefix(strings.TrimSpace(args[0]), ".")
	if newTLD == "" {
		return errors.New("tld cannot be empty")
	}
	wasSetup := config.IsSetup()
	old := config.GetTLD()
	if wasSetup && old != newTLD {
		_ = dns.New().Teardown(old)
	}
	if err := config.SetTLD(newTLD); err != nil {
		return err
	}
	if wasSetup {
		if err := apply(newTLD); err != nil {
			return err
		}
	}
	fmt.Printf("%s TLD is now .%s\n", green("✓"), newTLD)
	return nil
}

// List prints the currently mapped hosts.
func List() error {
	entries, err := os.ReadDir(config.SitesDir())
	if err != nil {
		fmt.Println("no active projects")
		return nil
	}
	n := 0
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".caddy" {
			fmt.Printf("  http://%s\n", strings.TrimSuffix(e.Name(), ".caddy"))
			n++
		}
	}
	if n == 0 {
		fmt.Println("no active projects")
	}
	return nil
}

// Doctor reports the health of the setup.
func Doctor() error {
	fmt.Printf("TLD          .%s\n", config.GetTLD())
	if config.IsSetup() {
		fmt.Println(green("✓") + " setup present")
	} else {
		fmt.Println("• not set up — run localtld setup")
	}
	fmt.Printf("Resolver     %s\n", dns.New().Name())
	return nil
}

// Uninstall reverts DNS changes, removes services, and clears local config.
func Uninstall() error {
	svc := service.New()
	_ = svc.Uninstall("caddy")
	_ = svc.Uninstall("dns")
	if err := dns.New().Teardown(config.GetTLD()); err != nil {
		return err
	}
	_ = os.RemoveAll(config.Dir())
	fmt.Println(green("✓") + " reverted DNS changes, removed services and config")
	return nil
}
