package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/abdullahharunozturk/localtld/internal/safety"
)

// confirmTLD blocks unless a real-TLD choice is explicitly confirmed. The scope
// of the warning matches the actual blast radius:
//
//	single label (com)     → hijacks ALL *.com  → hard "I UNDERSTAND"
//	multi label (saka.com) → only *.saka.com    → soft [y/N]
//	not a real TLD         → safe, no prompt
func confirmTLD(tld string) error {
	labels := strings.Split(tld, ".")
	ext := labels[len(labels)-1]
	if !safety.IsRealTLD(ext) {
		return nil
	}
	if len(labels) == 1 {
		return hardConfirm(ext)
	}
	return softConfirm(tld, ext)
}

func hardConfirm(ext string) error {
	if !isTTY(os.Stdin) {
		return fmt.Errorf("refusing to hijack all *.%s (a real TLD) without confirmation", ext)
	}
	fmt.Fprintf(os.Stderr, "\n%s .%s is a REAL TLD.\n", paint("41;1", " DANGER "), ext)
	fmt.Fprintf(os.Stderr, "  Every *.%s domain would resolve to 127.0.0.1 on this machine —\n", ext)
	fmt.Fprintf(os.Stderr, "  you'd lose access to all real .%s sites.\n", ext)
	fmt.Fprint(os.Stderr, `  Type "I UNDERSTAND" to proceed: `)
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimSpace(line) != "I UNDERSTAND" {
		return errors.New("aborted")
	}
	return nil
}

func softConfirm(tld, ext string) error {
	if !isTTY(os.Stdin) {
		return fmt.Errorf("refusing to route *.%s (under real TLD .%s) without confirmation", tld, ext)
	}
	fmt.Fprintf(os.Stderr, "\n%s .%s is a real TLD, but only *.%s points to 127.0.0.1.\n",
		paint("43;1", " NOTE "), ext, tld)
	fmt.Fprintf(os.Stderr, "  Other .%s sites are unaffected.\n", ext)
	fmt.Fprint(os.Stderr, "  Proceed? [y/N] ")
	var ans string
	_, _ = fmt.Fscanln(os.Stdin, &ans)
	switch strings.ToLower(strings.TrimSpace(ans)) {
	case "y", "yes":
		return nil
	default:
		return errors.New("aborted")
	}
}
