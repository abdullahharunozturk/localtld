package resolver

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func query(t *testing.T, addr, name string, qtype uint16) *dns.Msg {
	t.Helper()
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(name), qtype)
	r, _, err := c.Exchange(m, addr)
	if err != nil {
		t.Fatalf("exchange %s: %v", name, err)
	}
	return r
}

func TestServeAnswersTLD(t *testing.T) {
	addr := "127.0.0.1:15353"
	go func() { _ = Serve(addr, "localtld") }()
	time.Sleep(250 * time.Millisecond)

	// Wildcard subdomain → 127.0.0.1
	r := query(t, addr, "panel.aaron.localtld", dns.TypeA)
	if len(r.Answer) != 1 {
		t.Fatalf("want 1 A answer, got %d", len(r.Answer))
	}
	if a, ok := r.Answer[0].(*dns.A); !ok || a.A.String() != "127.0.0.1" {
		t.Fatalf("want 127.0.0.1, got %v", r.Answer[0])
	}

	// AAAA → ::1
	r = query(t, addr, "x.localtld", dns.TypeAAAA)
	if len(r.Answer) != 1 {
		t.Fatalf("want 1 AAAA answer, got %d", len(r.Answer))
	}

	// Foreign domain → no answer (we don't hijack the internet)
	r = query(t, addr, "example.com", dns.TypeA)
	if len(r.Answer) != 0 {
		t.Fatalf("want 0 answers for example.com, got %d", len(r.Answer))
	}
}
