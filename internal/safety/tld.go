// Package safety guards against hijacking a real TLD.
package safety

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdullahharunozturk/localtld/internal/config"
)

// fastList is an offline set of common real TLDs for an instant check with no
// network. The full IANA list (downloaded + cached) covers the rest.
var fastList = map[string]bool{
	"com": true, "net": true, "org": true, "io": true, "dev": true, "app": true,
	"co": true, "tech": true, "sh": true, "me": true, "ai": true, "info": true,
	"biz": true, "xyz": true, "online": true, "site": true, "store": true,
	"cloud": true, "tv": true, "uk": true, "us": true, "de": true, "fr": true,
	"nl": true, "es": true, "it": true,
}

// IsRealTLD reports whether ext (a single label) is a real IANA TLD.
func IsRealTLD(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	if ext == "" {
		return false
	}
	if fastList[ext] {
		return true
	}
	return ianaList()[ext]
}

const ianaURL = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"

func cachePath() string { return filepath.Join(config.Dir(), "iana-tlds.txt") }

// ianaList returns the IANA TLD set from cache, downloading it once if absent.
// Best-effort: on any failure it returns nil and only fastList applies.
func ianaList() map[string]bool {
	b, err := os.ReadFile(cachePath())
	if err != nil {
		if b, err = download(); err != nil {
			return nil
		}
	}
	m := map[string]bool{}
	sc := bufio.NewScanner(strings.NewReader(string(b)))
	for sc.Scan() {
		line := strings.TrimSpace(strings.ToLower(sc.Text()))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m[line] = true
	}
	return m
}

func download() ([]byte, error) {
	c := &http.Client{Timeout: 5 * time.Second}
	resp, err := c.Get(ianaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = os.MkdirAll(config.Dir(), 0o755)
	_ = os.WriteFile(cachePath(), b, 0o644)
	return b, nil
}
