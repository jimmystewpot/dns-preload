package dns

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Resolver interface {
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupNS(ctx context.Context, host string) ([]*net.NS, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
}

type resolver struct {
	client *net.Resolver
}

func NewResolver(nameserver string, timeout time.Duration) *resolver {
	return &resolver{
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
	}
}

// LookupCNAME returns the net.Resolver LookupCNAME
func (r *resolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	return r.client.LookupCNAME(ctx, host)
}

// LookupIPAddr returns the net.Resolver LookupIPAddr
func (r *resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.client.LookupIPAddr(ctx, host)
}

// LookupMX returns the net.Resolver LookupMX
func (r *resolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	return r.client.LookupMX(ctx, host)
}

// LookupTXT returns the net.Resolver LookupTXT
func (r *resolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	return r.client.LookupTXT(ctx, host)
}

// LookupNS returns the net.Resolver LookupNS
func (r *resolver) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
	return r.client.LookupNS(ctx, host)
}

// NewMockResolver returns the mock resolver.
func NewMockResolver() *mockresolver {
	return &mockresolver{}
}

// mock the above resolver interface
type mockresolver struct{}

func (m *mockresolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	switch host {
	case "www.foo.bar":
		return "foo.bar", nil
	case "www.bar.foo":
		return "", fmt.Errorf("no cname found")
	}
	return "", nil
}

func (m *mockresolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	switch host {
	case "foo.bar":
		ip1 := net.ParseIP("8.8.4.4")
		return []net.IPAddr{
			{
				IP: ip1,
			},
		}, nil
	case "mx0.foo.bar":
		ip1 := net.ParseIP("8.8.4.4")
		return []net.IPAddr{
			{
				IP: ip1,
			},
		}, nil
	case "mx1.foo.bar":
		return []net.IPAddr{}, fmt.Errorf("error")
	case "ns1.foo.bar":
		ip1 := net.ParseIP("8.8.4.4")
		return []net.IPAddr{
			{
				IP: ip1,
			},
		}, nil
	case "ns2.foo.bar":
		ip2 := net.ParseIP("8.8.8.8")
		return []net.IPAddr{
			{
				IP: ip2,
			},
		}, nil
	case "bar.foo":
		return []net.IPAddr{}, fmt.Errorf("error")

	}
	return []net.IPAddr{}, nil
}

func (m *mockresolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	switch host {
	case "foo.bar":
		return []*net.MX{
			{
				Host: "mx0.foo.bar",
				Pref: uint16(10),
			},
			{
				Host: "mx1.foo.bar",
				Pref: uint16(10),
			},
		}, nil
	case "bar.foo":
		return []*net.MX{}, fmt.Errorf("error")

	}
	return []*net.MX{}, nil
}

func (m *mockresolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	switch host {
	case "foo.bar":
		return []string{
			"v=spf1 -all",
		}, nil
	case "bar.foo":
		return []string{}, fmt.Errorf("error")
	}
	return []string{}, nil
}

func (m *mockresolver) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
	switch host {
	case "foo.bar":
		return []*net.NS{
			{
				Host: "ns1.foo.bar",
			},
			{
				Host: "ns2.foo.bar",
			},
		}, nil
	case "bar.foo":
		return []*net.NS{}, fmt.Errorf("error")
	}
	return []*net.NS{}, nil
}
