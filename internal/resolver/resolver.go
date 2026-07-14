// Package resolver is localtld's built-in DNS server. It answers every query
// under the configured TLD with the loopback address, replacing a per-OS
// dnsmasq dependency (dnsmasq isn't native on Windows). The OS resolver wiring
// (/etc/resolver, systemd-resolved, NRPT) points *.tld at 127.0.0.1:53 → here.
package resolver

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

// Serve runs the DNS server on addr (e.g. "127.0.0.1:53"), answering *.tld and
// tld itself with loopback. It blocks until a listener fails.
func Serve(addr, tld string) error {
	h := &handler{suffix: "." + strings.ToLower(dns.Fqdn(tld)), apex: strings.ToLower(dns.Fqdn(tld))}
	mux := dns.NewServeMux()
	mux.HandleFunc(".", h.handle)

	udp := &dns.Server{Addr: addr, Net: "udp", Handler: mux}
	tcp := &dns.Server{Addr: addr, Net: "tcp", Handler: mux}

	errc := make(chan error, 2)
	go func() { errc <- udp.ListenAndServe() }()
	go func() { errc <- tcp.ListenAndServe() }()
	return <-errc
}

type handler struct {
	suffix string // ".<tld>."
	apex   string // "<tld>."
}

func (h *handler) owns(name string) bool {
	name = strings.ToLower(name)
	return name == h.apex || strings.HasSuffix(name, h.suffix)
}

func (h *handler) handle(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		if !h.owns(q.Name) {
			continue
		}
		switch q.Qtype {
		case dns.TypeA:
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
				A:   net.IPv4(127, 0, 0, 1),
			})
		case dns.TypeAAAA:
			m.Answer = append(m.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 0},
				AAAA: net.IPv6loopback,
			})
		}
	}
	_ = w.WriteMsg(m)
}
