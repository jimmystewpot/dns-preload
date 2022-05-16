package dns

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"
)

var abc123 testing.T

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
			m := &mockresolver{}
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

func TestMockresolverLookupIPAddr(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		m       *mockresolver
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
					IP: net.ParseIP(googlePubDns1),
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
					IP: net.ParseIP(googlePubDns2),
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
			m := &mockresolver{}
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
		m       *mockresolver
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
			m := &mockresolver{}
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
		m       *mockresolver
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
			m := &mockresolver{}
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
		m       *mockresolver
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

func TestNewResolver(t *testing.T) {
	var testDomain string = "google.com"
	type args struct {
		nameserver string
		timeout    time.Duration
	}
	tests := []struct {
		name string
		args args
		want *resolver
	}{
		{
			name: "test",
			args: args{
				nameserver: "192.168.1.252:53",
				timeout:    time.Duration(5 * time.Second),
			},
			want: &resolver{client: &net.Resolver{}},
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
		})
	}
}
