// Package config holds localtld's machine-level settings and project-label lookup.
//
//   - TLD is machine-level  (~/.config/localtld/config → "tld=...")
//   - label is project-level (package.json → "localtld": "panel.aaron")
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	App        = "localtld"
	DefaultTLD = "localtld"
)

// Dir is the machine-level config directory.
// Unix: $XDG_CONFIG_HOME/localtld or ~/.config/localtld (matches the bash tool).
// Windows: %AppData%\localtld.
func Dir() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, App)
	}
	if runtime.GOOS == "windows" {
		if d, err := os.UserConfigDir(); err == nil {
			return filepath.Join(d, App)
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", App)
}

func File() string          { return filepath.Join(Dir(), "config") }
func SitesDir() string      { return filepath.Join(Dir(), "sites") }
func CaddyfilePath() string { return filepath.Join(Dir(), "Caddyfile") }

// IsSetup reports whether `localtld setup` has run on this machine.
func IsSetup() bool {
	_, err := os.Stat(File())
	return err == nil
}

// GetTLD returns the machine TLD, defaulting to DefaultTLD when unset.
func GetTLD() string {
	b, err := os.ReadFile(File())
	if err != nil {
		return DefaultTLD
	}
	for _, line := range strings.Split(string(b), "\n") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), "tld="); ok {
			if v = strings.TrimSpace(v); v != "" {
				return v
			}
		}
	}
	return DefaultTLD
}

// SetTLD persists the machine TLD.
func SetTLD(tld string) error {
	if err := os.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	return os.WriteFile(File(), []byte("tld="+tld+"\n"), 0o644)
}

// ReadLabel returns the "localtld" field from package.json in dir, or "".
func ReadLabel(dir string) string {
	b, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return ""
	}
	var pkg struct {
		Localtld string `json:"localtld"`
	}
	if json.Unmarshal(b, &pkg) != nil {
		return ""
	}
	return strings.TrimSpace(pkg.Localtld)
}
