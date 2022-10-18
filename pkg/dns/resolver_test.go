package dns

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
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

// NewMockResolver returns the mock resolver.
func NewMockResolver() *Mockresolver {
	return &Mockresolver{}
}

// Mock the above resolver interface
type Mockresolver struct{}

func (m *Mockresolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	switch host {
	case "www.foo.bar":
		return testDomainNoErr, nil
	case "www.bar.foo":
		return "", fmt.Errorf(nxDomainErr, host)
	}
	return "", fmt.Errorf(nxDomainErr, host)
}

func (m *Mockresolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	switch host {
	case testDomainNoErr, testDomainMX0, testDomainNS1:
		ip1 := net.ParseIP(googlePubDNS1)
		return []net.IPAddr{
			{
				IP: ip1,
			},
		}, nil
	case "ns2.foo.bar":
		ip2 := net.ParseIP(googlePubDNS2)
		return []net.IPAddr{
			{
				IP: ip2,
			},
		}, nil
	case "dns.oranged.to":
		ip1 := net.ParseIP(googlePubDNS1)
		ip2 := net.ParseIP(googlePubDNS2)
		return []net.IPAddr{
			{
				IP: ip1,
			},
			{
				IP: ip2,
			},
		}, nil
	case testDomainWithErr:
		return []net.IPAddr{}, fmt.Errorf(nxDomainErr, host)
	}
	return []net.IPAddr{}, fmt.Errorf(nxDomainErr, host)
}

func (m *Mockresolver) LookupAddr(ctx context.Context, addr string) ([]string, error) {
	if addr != googleIpv6 {
		return []string{}, fmt.Errorf("%s ptr not found", addr)
	}
	return []string{"ipv6.google.com"}, nil
}

//nolint:gocritic // uses switch to expand on test cases in the future.
func (m *Mockresolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	switch host {
	case testDomainNoErr:
		return []*net.MX{
			{
				Host: testDomainMX0,
				Pref: 10,
			},
			{
				Host: testDomainMX1,
				Pref: 10,
			},
		}, nil
	}
	return []*net.MX{}, fmt.Errorf(nxDomainErr, host)
}

//nolint:gocritic // uses switch to expand on test cases in the future.
func (m *Mockresolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	switch host {
	case testDomainNoErr:
		return []string{
			"v=spf1 -all",
		}, nil
	}
	return []string{}, fmt.Errorf(nxDomainErr, host)
}

//nolint:gocritic // uses switch to expand on test cases in the future.
func (m *Mockresolver) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
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

func TestResolverLookupAll(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid cname",
			args: args{
				ctx:  context.Background(),
				host: "www.foo.bar",
			},
			wantErr: false,
			want:    "foo.bar",
		},
		{
			name: "invalid cname",
			args: args{
				ctx:  context.Background(),
				host: "www.bar.foo",
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "nxdomain",
			args: args{
				ctx:  context.Background(),
				host: "www.xxx.yyy",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mockresolver{}
			cname, err := m.LookupCNAME(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockresolver.LookupCNAME() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cname != tt.want {
				t.Errorf("mockresolver.LookupCNAME() = %v, want %v", cname, tt.want)
			}
		})
	}
}

func TestMockresolverLookupAddr(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		m       *Mockresolver
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid a lookup",
			args: args{
				ctx:  context.Background(),
				host: testPtrNoErr,
			},
			want:    []string{"ipv6.google.com"},
			wantErr: false,
		},
		{
			name: "invalid a lookup",
			args: args{
				ctx:  context.Background(),
				host: testDomainWithErr,
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mockresolver{}
			got, err := m.LookupAddr(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockresolver.LookupAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockresolver.LookupAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockresolverLookupIPAddr(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		m       *Mockresolver
		args    args
		want    []net.IPAddr
		wantErr bool
	}{
		{
			name: "valid a lookup",
			args: args{
				ctx:  context.Background(),
				host: testDomainNoErr,
			},
			want: []net.IPAddr{
				{
					IP: net.ParseIP(googlePubDNS1),
				},
			},
			wantErr: false,
		},
		{
			name: "valid a lookup 2",
			args: args{
				ctx:  context.Background(),
				host: "ns2.foo.bar",
			},
			want: []net.IPAddr{
				{
					IP: net.ParseIP(googlePubDNS2),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid a lookup",
			args: args{
				ctx:  context.Background(),
				host: testDomainWithErr,
			},
			want:    []net.IPAddr{},
			wantErr: true,
		},
		{
			name: "nxdomain",
			args: args{
				ctx:  context.Background(),
				host: "yyy.xxx.www",
			},
			want:    []net.IPAddr{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mockresolver{}
			got, err := m.LookupIPAddr(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockresolver.LookupIPAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockresolver.LookupIPAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockresolverLookupMX(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		m       *Mockresolver
		args    args
		want    []*net.MX
		wantErr bool
	}{
		{
			name: "A valid MX",
			args: args{
				ctx:  context.Background(),
				host: testDomainNoErr,
			},
			want: []*net.MX{
				{
					Host: testDomainMX0,
					Pref: 10,
				},
				{
					Host: testDomainMX1,
					Pref: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "An invalid MX",
			args: args{
				ctx:  context.Background(),
				host: testDomainWithErr,
			},
			want:    []*net.MX{},
			wantErr: true,
		},
		{
			name: "nxdomain",
			args: args{
				ctx:  context.Background(),
				host: "aaa.111.ccc",
			},
			want:    []*net.MX{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mockresolver{}
			got, err := m.LookupMX(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockresolver.LookupMX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockresolver.LookupMX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockresolverLookupTXT(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		m       *Mockresolver
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "a valid txt record",
			args: args{
				ctx:  context.Background(),
				host: testDomainNoErr,
			},
			want: []string{
				"v=spf1 -all",
			},
			wantErr: false,
		},
		{
			name: "a invalid txt record",
			args: args{
				ctx:  context.Background(),
				host: testDomainWithErr,
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "nxdomain",
			args: args{
				ctx:  context.Background(),
				host: "xxx.101-4-1",
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mockresolver{}
			got, err := m.LookupTXT(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockresolver.LookupTXT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockresolver.LookupTXT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockresolverLookupNS(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		m       *Mockresolver
		args    args
		want    []*net.NS
		wantErr bool
	}{
		{
			name: "valid ns lookup",
			args: args{
				ctx:  context.Background(),
				host: testDomainNoErr,
			},
			want: []*net.NS{
				{
					Host: "ns1.foo.bar",
				},
				{
					Host: "ns2.foo.bar",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid ns lookup",
			args: args{
				ctx:  context.Background(),
				host: testDomainWithErr,
			},
			want:    []*net.NS{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockResolver()
			got, err := m.LookupNS(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockresolver.LookupNS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockresolver.LookupNS() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:errcheck,gocritic // errcheck here is just cycling through various tests
func TestNewResolver(t *testing.T) {
	var testDomain = "google.com"
	type args struct {
		nameserver string
		timeout    time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Resolver
	}{
		{
			name: "test",
			args: args{
				nameserver: "192.168.1.252:53",
				timeout:    5 * time.Second,
			},
			want: &Resolver{client: &net.Resolver{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewResolver(tt.args.nameserver, tt.args.timeout)
			if !reflect.DeepEqual(resolver, resolver) {
				t.Errorf("NewResolver() = %v, want %v", resolver, tt.want)
			}
			resolver.LookupIPAddr(context.Background(), testDomain)
			resolver.LookupMX(context.Background(), testDomain)
			resolver.LookupTXT(context.Background(), testDomain)
			resolver.LookupNS(context.Background(), testDomain)
			resolver.LookupAddr(context.Background(), testPtrNoErr)
			resolver.LookupCNAME(context.Background(), testDomain)
		})
	}
}
