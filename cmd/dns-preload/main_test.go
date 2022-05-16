package main

import (
	"context"
	"testing"
	"time"

	"github.com/jimmystewpot/dns-preload/pkg/dns"
)

const (
	testDomainNoErr   string = "foo.bar"
	testDomainWithErr string = "bar.foo"
)

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
		resolver   dns.Resolver
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
				resolver: dns.NewMockResolver(),
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
				resolver: dns.NewMockResolver(),
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
				Delay:      tt.fields.Delay,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.Hosts(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.Hosts() error = %v, wantErr %v", err, tt.wantErr)
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
		resolver   dns.Resolver
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
				resolver: dns.NewMockResolver(),
				Full:     false,
				Quiet:    false,
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
				resolver: dns.NewMockResolver(),
				Full:     true,
				Quiet:    false,
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
				resolver: dns.NewMockResolver(),
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
				Delay:      tt.fields.Delay,
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
		resolver   dns.Resolver
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
				resolver: dns.NewMockResolver(),
				Full:     true,
				Quiet:    false,
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
				resolver: dns.NewMockResolver(),
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
				Delay:      tt.fields.Delay,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
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
		resolver   dns.Resolver
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
				resolver: dns.NewMockResolver(),
				Full:     true,
				Quiet:    false,
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
				resolver: dns.NewMockResolver(),
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
				Delay:      tt.fields.Delay,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
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
		resolver   dns.Resolver
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
				resolver: dns.NewMockResolver(),
				Full:     true,
				Quiet:    false,
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
				resolver: dns.NewMockResolver(),
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
				Delay:      tt.fields.Delay,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.CNAME(tt.args.ctx, tt.args.hosts); (err != nil) != tt.wantErr {
				t.Errorf("Preload.CNAME() error = %v, wantErr %v", err, tt.wantErr)
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
		Quiet      bool
		Full       bool
		Debug      bool
		Timeout    time.Duration
		Delay      time.Duration
		resolver   dns.Resolver
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
			name:    "err test",
			fields:  fields{},
			args:    args{},
			wantErr: true,
		},
		{
			name: "good config test",
			fields: fields{
				ConfigFile: "../../pkg/confighandlers/test_data/complete_config_sample.yaml",
				Server:     "192.168.1.252",
				Port:       "53",
			},
			args: args{
				cmd: "cname",
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
				Delay:      tt.fields.Delay,
				resolver:   tt.fields.resolver,
				nameserver: tt.fields.nameserver,
			}
			if err := p.Run(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Preload.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
