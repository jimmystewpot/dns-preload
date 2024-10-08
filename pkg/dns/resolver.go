package dns

import (
	"context"
	"net"
	"time"
)

// CustomResolver Interface can be reimplemented very easily for mocking
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
