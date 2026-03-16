package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

// CustomResolver Interface can be reimplemented very easily for mocking
type CustomResolver interface {
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupAddr(ctx context.Context, addr string) ([]string, error)
	LookupNS(ctx context.Context, host string) ([]*net.NS, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
	LookupCNAMEWithDNSSEC(ctx context.Context, host string) (string, error)
	LookupIPAddrWithDNSSEC(ctx context.Context, host string) ([]net.IPAddr, error)
}

type Resolver struct {
	client     *net.Resolver
	nameserver string
	timeout    time.Duration
}

// NewResolver creates a custom resolver where the DNS servers are pinned.
func NewResolver(nameserver string, timeout time.Duration) *Resolver {
	//nolint:revive // address is a returned function, it gets set by the caller.
	return &Resolver{
		client: &net.Resolver{
			PreferGo:     true,
			StrictErrors: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: timeout,
				}
				return d.DialContext(ctx, network, nameserver)
			},
		},
		nameserver: nameserver,
		timeout:    timeout,
	}
}

// LookupCNAME returns the net.Resolver LookupCNAME
func (r *Resolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	return r.client.LookupCNAME(ctx, host)
}

// LookupAddr returns the net.Resolver LookupAddr
func (r *Resolver) LookupAddr(ctx context.Context, host string) ([]string, error) {
	return r.client.LookupAddr(ctx, host)
}

// LookupIPAddr returns the net.Resolver LookupIPAddr
func (r *Resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.client.LookupIPAddr(ctx, host)
}

// LookupMX returns the net.Resolver LookupMX
func (r *Resolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	return r.client.LookupMX(ctx, host)
}

// LookupTXT returns the net.Resolver LookupTXT
func (r *Resolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	return r.client.LookupTXT(ctx, host)
}

// LookupNS returns the net.Resolver LookupNS
func (r *Resolver) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
	return r.client.LookupNS(ctx, host)
}

// LookupCNAMEWithDNSSEC performs CNAME lookup with DNSSEC validation
func (r *Resolver) LookupCNAMEWithDNSSEC(ctx context.Context, host string) (string, error) {
	c := new(dns.Client)
	c.Timeout = r.timeout
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeCNAME)
	m.SetEdns0(4096, true) // Enable DNSSEC
	resp, _, err := c.Exchange(m, r.nameserver)
	if err != nil {
		return "", err
	}
	if resp.Rcode != dns.RcodeSuccess {
		return "", fmt.Errorf("DNS query failed with rcode %d", resp.Rcode)
	}
	// Basic DNSSEC validation (check for AD bit or RRSIG)
	if !resp.AuthenticatedData {
		return "", fmt.Errorf("DNSSEC validation failed: not authenticated")
	}
	for _, ans := range resp.Answer {
		if cname, ok := ans.(*dns.CNAME); ok {
			return cname.Target, nil
		}
	}
	return "", fmt.Errorf("no CNAME record found")
}

// LookupIPAddrWithDNSSEC performs IP address lookup with DNSSEC validation
func (r *Resolver) LookupIPAddrWithDNSSEC(ctx context.Context, host string) ([]net.IPAddr, error) {
	c := new(dns.Client)
	c.Timeout = r.timeout
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)
	m.SetEdns0(4096, true)
	resp, _, err := c.Exchange(m, r.nameserver)
	if err != nil {
		return nil, err
	}
	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS query failed with rcode %d", resp.Rcode)
	}
	if !resp.AuthenticatedData {
		return nil, fmt.Errorf("DNSSEC validation failed: not authenticated")
	}
	var addrs []net.IPAddr
	for _, ans := range resp.Answer {
		if a, ok := ans.(*dns.A); ok {
			addrs = append(addrs, net.IPAddr{IP: a.A})
		}
	}
	// Also query AAAA
	m.SetQuestion(dns.Fqdn(host), dns.TypeAAAA)
	resp, _, err = c.Exchange(m, r.nameserver)
	if err != nil {
		return addrs, nil // Return A records even if AAAA fails
	}
	if resp.Rcode == dns.RcodeSuccess && resp.AuthenticatedData {
		for _, ans := range resp.Answer {
			if aaaa, ok := ans.(*dns.AAAA); ok {
				addrs = append(addrs, net.IPAddr{IP: aaaa.AAAA})
			}
		}
	}
	return addrs, nil
}
