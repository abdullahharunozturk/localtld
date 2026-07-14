package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abdullahharunozturk/localtld/internal/config"
	"github.com/abdullahharunozturk/localtld/internal/dns"
	"github.com/abdullahharunozturk/localtld/internal/proxy"
)

// Setup configures DNS + Caddy for the current machine TLD.
func Setup() error {
	tld := config.GetTLD()
	fmt.Println(bold("localtld setup"))
	fmt.Printf("  TLD: .%s\n", tld)

	if err := config.SetTLD(tld); err != nil {
		return err
	}
	if err := proxy.WriteBase(); err != nil {
		return err
	}
	p := dns.New()
	fmt.Printf("  DNS: %s\n", p.Name())
	if err := p.Setup(tld); err != nil {
		return fmt.Errorf("dns setup: %w", err)
	}
	fmt.Println(green("✓ setup complete"))
	fmt.Println(dim(`  Add "localtld": "myapp" to package.json, then: localtld run -- <dev cmd>`))
	return nil
}

// TLD changes the machine-wide TLD, re-wiring DNS if already set up.
func TLD(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: localtld tld <new-tld>   (e.g. localtld tld test)")
	}
	newTLD := strings.TrimPrefix(strings.TrimSpace(args[0]), ".")
	if newTLD == "" {
		return errors.New("tld cannot be empty")
	}
	old := config.GetTLD()
	if config.IsSetup() && old != newTLD {
		_ = dns.New().Teardown(old)
	}
	if err := config.SetTLD(newTLD); err != nil {
		return err
	}
	if config.IsSetup() {
		if err := dns.New().Setup(newTLD); err != nil {
			return err
		}
	}
	fmt.Printf("%s TLD is now .%s\n", green("✓"), newTLD)
	return nil
}

// List prints the currently mapped hosts.
func List() error {
	entries, err := os.ReadDir(config.SitesDir())
	if err != nil || len(entries) == 0 {
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
	fmt.Printf("DNS provider %s\n", dns.New().Name())
	return nil
}

// Uninstall reverts DNS changes and removes local config.
func Uninstall() error {
	if err := dns.New().Teardown(config.GetTLD()); err != nil {
		return err
	}
	_ = os.RemoveAll(config.Dir())
	fmt.Println(green("✓") + " reverted DNS changes and removed config")
	return nil
}
