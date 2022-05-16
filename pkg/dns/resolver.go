package dns

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	googlePubDns1     string = "8.8.4.4"
	googlePubDns2     string = "8.8.8.8"
	testDomainNoErr   string = "foo.bar"
	testDomainWithErr string = "bar.foo"
	testDomainMX0     string = "mx0.foo.bar"
	testDomainMX1     string = "mx1.foo.bar"
	testDomainNS1     string = "ns1.foo.bar"
	nxDomainErr       string = "nxdomain %s"
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
		return testDomainNoErr, nil
	case "www.bar.foo":
		return "", fmt.Errorf(nxDomainErr, host)
	}
	return "", fmt.Errorf(nxDomainErr, host)
}

func (m *mockresolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	switch host {
	case testDomainNoErr, testDomainMX0, testDomainNS1:
		ip1 := net.ParseIP(googlePubDns1)
		return []net.IPAddr{
			{
				IP: ip1,
			},
		}, nil
	case "ns2.foo.bar":
		ip2 := net.ParseIP(googlePubDns2)
		return []net.IPAddr{
			{
				IP: ip2,
			},
		}, nil
	case testDomainWithErr:
		return []net.IPAddr{}, fmt.Errorf(nxDomainErr, host)
	}
	return []net.IPAddr{}, fmt.Errorf(nxDomainErr, host)
}

func (m *mockresolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	switch host {
	case testDomainNoErr:
		return []*net.MX{
			{
				Host: testDomainMX0,
				Pref: uint16(10),
			},
			{
				Host: testDomainMX1,
				Pref: uint16(10),
			},
		}, nil
	}
	return []*net.MX{}, fmt.Errorf(nxDomainErr, host)
}

func (m *mockresolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	switch host {
	case testDomainNoErr:
		return []string{
			"v=spf1 -all",
		}, nil
	}
	return []string{}, fmt.Errorf(nxDomainErr, host)
}

func (m *mockresolver) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
	switch host {
	case testDomainNoErr:
		return []*net.NS{
			{
				Host: "ns1.foo.bar",
			},
			{
				Host: "ns2.foo.bar",
			},
		}, nil
	}
	return []*net.NS{}, fmt.Errorf(nxDomainErr, host)
}
