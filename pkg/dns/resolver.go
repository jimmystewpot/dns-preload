package dns

import (
	"context"
	"net"
	"time"
)

const (
	googlePubDNS1     string = "8.8.4.4"
	googlePubDNS2     string = "8.8.8.8"
	googleIpv6        string = "2404:6800:4006:804::200e"
	testDomainNoErr   string = "foo.bar"
	testDomainWithErr string = "bar.foo"
	testDomainMX0     string = "mx0.foo.bar"
	testDomainMX1     string = "mx1.foo.bar"
	testDomainNS1     string = "ns1.foo.bar"
	testPtrNoErr      string = "2404:6800:4006:804::200e"
	nxDomainErr       string = "nxdomain %s"
)

type CustomResolver interface {
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupAddr(ctx context.Context, addr string) ([]string, error)
	LookupNS(ctx context.Context, host string) ([]*net.NS, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
}

type Resolver struct {
	client *net.Resolver
}

func NewResolver(nameserver string, timeout time.Duration) *Resolver {
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
