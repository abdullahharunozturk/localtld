// Package netutil holds small networking helpers.
package netutil

import "net"

// FreePort asks the OS for an unused TCP port on the loopback interface.
// The kernel guarantees uniqueness at allocation time; the dev server binds
// it a moment later via the PORT env we pass.
func FreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
