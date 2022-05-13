package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/jimmystewpot/dns-preload/pkg/confighandlers"
)

const (
	// these const strings are used to store the DNS query types for reuse.
	aTypeStr       string = "A, AAAA"
	cnameTypeStr   string = "CNAME"
	mxTypeStr      string = "MX"
	nsTypeStr      string = "NS"
	txtTypeStr     string = "TXT"
	preloadMessage string = "Preloading Nameserver: %s with query type: %s for domains: %s\n"
)

var (
	cli struct {
		All   Preload `cmd:""`
		Cname Preload `cmd:""`
		Hosts Preload `cmd:""`
		Mx    Preload `cmd:""`
		Ns    Preload `cmd:""`
		Txt   Preload `cmd:""`
	}

	start time.Time
)

type Preload struct {
	ConfigFile string        `required:"" help:"The configuration file to read the domain list to query from"`
	Server     string        `default:"localhost" help:"The server to query to seed the domain list into"`
	Port       string        `default:"53" help:"The port to query for on the DNS server"`
	Workers    uint8         `default:"5" help:"The number of concurrent goroutines used to query the DNS server"`
	Quiet      bool          `default:"false" help:"Suppress the preload response output to console"`
	Full       bool          `default:"true" help:"For record types that return a Hostname ensure that these are resolved"`
	Timeout    time.Duration `default:"30s" help:"The timeout for DNS queries to succeed"`
	Delay      time.Duration `default:"0s" help:"How long to wait until the queries are executed"`
	resolver   *net.Resolver
	nameserver string
}

func (p *Preload) Run(cmd string) error {
	cfg, err := confighandlers.LoadConfigFromFile(&p.ConfigFile)
	if err != nil {
		return err
	}

	p.nameserver = net.JoinHostPort(p.Server, p.Port)
	p.resolver = &net.Resolver{
		PreferGo:     true,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: p.Timeout,
			}
			return d.DialContext(ctx, network, p.nameserver)
		},
	}
	ctx := context.Background()
	switch cmd {
	case "all":
		e := make([]error, 0)
		for _, query := range confighandlers.QueryTypes {
			err := p.Run(query)
			if err != nil {
				e = append(e, err)
			}
		}
		if len(e) != 0 {
			return fmt.Errorf("%s", e)
		}
	case "cname":
		if cfg.QueryType.CnameCount != 0 {
			fmt.Printf(preloadMessage, p.nameserver, cnameTypeStr, strings.Join(cfg.QueryType.Cname, ", "))
			return p.CNAME(ctx, cfg.QueryType.Cname)
		}
	case "hosts":
		if cfg.QueryType.HostsCount != 0 {
			fmt.Printf(preloadMessage, p.nameserver, aTypeStr, strings.Join(cfg.QueryType.Hosts, ", "))
			return p.Hosts(ctx, cfg.QueryType.Hosts)
		}
	case "mx":
		if cfg.QueryType.MXCount != 0 {
			fmt.Printf(preloadMessage, p.nameserver, mxTypeStr, strings.Join(cfg.QueryType.MX, ", "))
			return p.MX(ctx, cfg.QueryType.MX)
		}
	case "ns":
		if cfg.QueryType.NSCount != 0 {
			fmt.Printf(preloadMessage, p.nameserver, nsTypeStr, strings.Join(cfg.QueryType.NS, ", "))
			return p.NS(ctx, cfg.QueryType.NS)
		}
	case "txt":
		if cfg.QueryType.TXTCount != 0 {
			fmt.Printf(preloadMessage, p.nameserver, txtTypeStr, strings.Join(cfg.QueryType.TXT, ", "))
			return p.TXT(ctx, cfg.QueryType.TXT)
		}
	default:
		return fmt.Errorf("%s unknown command", cmd)
	}
	return nil
}

// Hosts preload the nameserver with IP addresses for a given list of hostnames.
func (p *Preload) Hosts(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		s := time.Now()
		result, err := p.resolver.LookupIPAddr(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], aTypeStr, time.Since(s), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// CNAME preload the nameserver with CNAME lookups for a given list of hostnames.
func (p *Preload) CNAME(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		s := time.Now()
		result, err := p.resolver.LookupCNAME(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], cnameTypeStr, time.Since(s), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// NS preloads the nameserver records for a given list of hostnames.
func (p *Preload) NS(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		s := time.Now()
		result, err := p.resolver.LookupNS(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], nsTypeStr, time.Since(s), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// MX preloads the nameserver with the MX records for a given list of hostnames.
func (p *Preload) MX(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		s := time.Now()
		result, err := p.resolver.LookupMX(ctx, hosts[i])
		if err != nil {
			return err
		}

		err = p.Printer(hosts[i], mxTypeStr, time.Since(s), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// TXT preloads the nameserver with the TXT records for a given list of hostnames.
func (p *Preload) TXT(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		s := time.Now()
		result, err := p.resolver.LookupTXT(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], txtTypeStr, time.Since(s), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// String provides output to the console for the results of the preloading.
func (p *Preload) Printer(hostname string, qtype string, duration time.Duration, results interface{}) error {
	// str is used to store the string conversions of the results.
	str := make([]string, 0)
	switch r := results.(type) {
	case string:
		str = append(str, results.(string))
	case []string:
		str = append(str, results.([]string)...)
	case []*net.MX:
		for _, mx := range results.([]*net.MX) {
			str = append(str, mx.Host)
			if p.Full {
				err := p.Hosts(context.Background(), str)
				if err != nil {
					return err
				}
			}
		}
	case []*net.NS:
		for _, ns := range results.([]*net.NS) {
			str = append(str, ns.Host)
			if p.Full {
				err := p.Hosts(context.Background(), str)
				if err != nil {
					return err
				}
			}
		}
	case []net.IPAddr:
		for _, ip := range results.([]net.IPAddr) {
			str = append(str, ip.IP.String())
		}
	default:
		return fmt.Errorf("error: got type %+v", r.(string))
	}
	if !p.Quiet {
		fmt.Printf("Preloaded %s type %s in %s to %+s\n", hostname, qtype, duration, strings.Join(str, ", "))
	}
	return nil
}

func main() {
	start = time.Now()
	cmd := kong.Parse(&cli,
		kong.Name(os.Args[0]),
		kong.Description("Preload a series of Domain Names into a DNS server from a yaml configuration"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	err := cmd.Run(cmd.Command())
	if err != nil {
		fmt.Println(err)
	}
	cmd.FatalIfErrorf(err)
	fmt.Printf("Preload completed in %s\n", time.Since(start))
}
