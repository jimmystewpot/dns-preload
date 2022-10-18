package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/jimmystewpot/dns-preload/pkg/confighandlers"
	"github.com/jimmystewpot/dns-preload/pkg/dns"
)

const (
	testDomainNoErr   string = "foo.bar"
	testDomainWithErr string = "bar.foo"
	googlePubDNS1     string = "8.8.4.4"
	googlePubDNS2     string = "8.8.8.8"
	googleIpv6        string = "2404:6800:4006:804::200e"
	testDNSServer     string = "9.9.9.9"
	testDNSServerPort string = "53"
	nxDomainErr       string = "nxdomain %s"
	testDomainMX0     string = "mx0.foo.bar"
	testDomainMX1     string = "mx1.foo.bar"
	testDomainNS1     string = "ns1.foo.bar"
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

func TestPreloadHosts(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		ctx   context.Context
		hosts []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Case Without Error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainNoErr},
			},
			wantErr: false,
		},
		{
			name: "Test Case With Error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainWithErr},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.Hosts(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.Hosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadPtr(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		ctx   context.Context
		hosts []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Case Without Error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{googleIpv6},
			},
			wantErr: false,
		},
		{
			name: "Test Case With Error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Quiet:      false,
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainNoErr},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.PTR(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.PTR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadMX(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		ctx   context.Context
		hosts []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "IN MX",
			fields: fields{
				resolver:   NewMockResolver(),
				Full:       false,
				Quiet:      false,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainNoErr},
			},
			wantErr: false,
		},
		{
			name: "IN MX full recursion with error",
			fields: fields{
				resolver:   NewMockResolver(),
				Full:       true,
				Quiet:      false,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx: ctx,
				// full recursive second lookup MX has a failure test case.
				hosts: []string{testDomainNoErr},
			},
			wantErr: true,
		},
		{
			name: "IN MX with Error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainWithErr},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.MX(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.Hosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadTXT(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		ctx   context.Context
		hosts []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "IN TXT",
			fields: fields{
				resolver:   NewMockResolver(),
				Full:       true,
				Quiet:      false,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainNoErr},
			},
			wantErr: false,
		},
		{
			name: "IN TXT error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainWithErr},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: net.JoinHostPort(tt.fields.Server, tt.fields.Port),
			}
			if err := p.TXT(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.TXT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadNS(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		ctx   context.Context
		hosts []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "IN NS",
			fields: fields{
				resolver:   NewMockResolver(),
				Full:       true,
				Quiet:      false,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainNoErr},
			},
			wantErr: false,
		},
		{
			name: "IN NS error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{testDomainWithErr},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: net.JoinHostPort(tt.fields.Server, tt.fields.Port),
			}
			if err := p.NS(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.NS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadCNAME(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		ctx   context.Context
		hosts []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "IN CNAME",
			fields: fields{
				resolver:   NewMockResolver(),
				Full:       true,
				Quiet:      false,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{"www.foo.bar"},
			},
			wantErr: false,
		},
		{
			name: "IN CNAME error",
			fields: fields{
				resolver:   NewMockResolver(),
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				ctx:   ctx,
				hosts: []string{"www.bar.foo"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.CNAME(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.CNAME() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadRunQueries(t *testing.T) {
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		nameserver string
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "good config test - cname",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				cmd: "cname",
			},
			wantErr: false,
		},
		{
			name: "good config test - hosts",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				cmd: "hosts",
			},
			wantErr: false,
		},
		{
			name: "good config test - txt",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				cmd: "txt",
			},
			wantErr: false,
		},
		{
			name: "good config test - mx",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				cmd: "mx",
			},
			wantErr: false,
		},
		{
			name: "good config test - ns",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Workers:    1,
			},
			args: args{
				cmd: "ns",
			},
			wantErr: false,
		},
		{
			name: "good config test - wrong cmd",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Debug:      true,
				Workers:    1,
			},
			args: args{
				cmd: "foo",
			},
			wantErr: true,
		},
		{
			name: "good config test - cname with no entries",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_no_cname_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Debug:      true,
				Workers:    1,
			},
			args: args{
				cmd: "cname",
			},
			wantErr: false,
		},
		{
			name: "good config test - ptr with entries",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_data_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Debug:      true,
				Workers:    1,
			},
			args: args{
				cmd: "ptr",
			},
			wantErr: false,
		},
		{
			name: "good config test - ptr with no entries",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/basic_test_no_cname_config.yaml",
				Server:     testDNSServer,
				Port:       testDNSServerPort,
				nameserver: net.JoinHostPort(testDNSServer, testDNSServerPort),
				Debug:      true,
				Workers:    1,
			},
			args: args{
				cmd: "ptr",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   NewMockResolver(),
				nameserver: tt.fields.nameserver,
			}
			cfg, _ := confighandlers.LoadConfigFromFile(&tt.fields.ConfigFile)
			if err := p.RunQueries(context.Background(), tt.args.cmd, cfg); (err != nil) != tt.wantErr {
				t.Errorf("Preload.RunQueries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadPrinter(t *testing.T) {
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   *Mockresolver
		nameserver string
	}
	type args struct {
		hostname string
		qtype    string
		duration time.Duration
		results  interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test other interface underlying data type",
			fields: fields{},
			args: args{
				hostname: testDomainNoErr,
				qtype:    "foo",
				duration: time.Second,
				results:  returnIntInterface(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.ResultsPrinter(tt.args.hostname, tt.args.qtype, tt.args.duration, tt.args.results); (err != nil) != tt.wantErr {
				t.Errorf("Preload.Printer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigRun(t *testing.T) {
	type fields struct {
		Quiet    bool
		Generate bool
		Validate struct{ ConfigFile string }
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Generate",
			fields: fields{
				Quiet: true,
			},
			args: args{
				cmd: "config generate",
			},
			wantErr: false,
		},
		{
			name: "Validate",
			fields: fields{
				Quiet: true,
			},
			args: args{
				cmd: "config validate",
			},
			wantErr: true,
		},
		{
			name: "error",
			fields: fields{
				Quiet: true,
			},
			args: args{
				cmd: "foo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Quiet: tt.fields.Quiet,
			}
			if err := c.Run(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Config.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreloadRun(t *testing.T) {
	type fields struct {
		ConfigFile string
		Server     string
		Port       string
		Workers    uint8
		Mute       bool
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		resolver   *dns.Resolver
		nameserver string
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "generate an error with good config",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/complete_config_sample.yaml",
				Workers:    1,
			},
			args: args{
				cmd: "does not exist",
			},
			wantErr: true,
		},
		{
			name: "generate an error with bad config",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/bad_configuration_sample.yaml",
				Workers:    1,
			},
			args: args{
				cmd: "does not exist",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Preload{
				ConfigFile: tt.fields.ConfigFile,
				Server:     tt.fields.Server,
				Port:       tt.fields.Port,
				Workers:    tt.fields.Workers,
				Mute:       tt.fields.Mute,
				Quiet:      tt.fields.Quiet,
				Full:       tt.fields.Full,
				Debug:      tt.fields.Debug,
				Timeout:    tt.fields.Timeout,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.Run(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Preload.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_completedPrinter(t *testing.T) {
	start := time.Now()
	type args struct {
		quiet bool
		t     time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test quiet mode",
			args: args{
				quiet: true,
				t:     start,
			},
			want: "",
		},
		{
			name: "test normal mode",
			args: args{
				quiet: false,
				t:     start,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completedPrinter(tt.args.quiet, tt.args.t)
		})
	}
}

func returnIntInterface() interface{} {
	x := []int{1, 2, 3, 4, 5}
	return x
}
