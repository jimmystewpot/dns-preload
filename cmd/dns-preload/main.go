package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/jimmystewpot/dns-preload/pkg/confighandlers"
)

const (
	// these const strings are used to store the DNS query types for reuse.
	aType     string = "A, AAAA"
	cnameType string = "CNAME"
	mxType    string = "MX"
	nsType    string = "NS"
	txtType   string = "TXT"
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
	// queryTypes is used to iterate through all of the commands when the all cmd is used.
	queryTypes = []string{"hosts", "cname", "mx", "ns", "txt"}
	start      time.Time
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
	wg         *sync.WaitGroup
}

func (p *Preload) Run(qtype string) error {
	cfg, err := confighandlers.LoadConfigFromFile(&p.ConfigFile)
	if err != nil {
		return err
	}
	start = time.Now()
	nameserver := net.JoinHostPort(p.Server, p.Port)
	p.resolver = &net.Resolver{
		PreferGo:     true,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: p.Timeout,
			}
			return d.DialContext(ctx, network, nameserver)
		},
	}
	ctx := context.Background()

	switch qtype {
	case "all":
		e := make([]error, 0)
		for _, query := range queryTypes {
			err := p.Run(query)
			if err != nil {
				e = append(e, err)
			}
		}
		if len(e) != 0 {
			return fmt.Errorf("%s", e)
		}
	case "cname":
		if !Starter(nameserver, qtype, cfg.QueryType.Cname) {
			return nil
		}
		return p.CNAME(ctx, cfg.QueryType.Cname)
	case "hosts":
		if !Starter(nameserver, aType, cfg.QueryType.Hosts) {
			return nil
		}
		return p.Hosts(ctx, cfg.QueryType.Hosts)
	case "mx":
		if !Starter(nameserver, qtype, cfg.QueryType.MX) {
			return nil
		}
		return p.MX(ctx, cfg.QueryType.MX)
	case "ns":
		if !Starter(nameserver, qtype, cfg.QueryType.NS) {
			return nil
		}
		return p.NS(ctx, cfg.QueryType.NS)
	case "txt":
		if !Starter(nameserver, qtype, cfg.QueryType.TXT) {
			return nil
		}
		return p.TXT(ctx, cfg.QueryType.TXT)
	default:
		return fmt.Errorf("%s incorrect command", qtype)
	}
	return nil
}

// Starter outputs what domains and query types we are about to preload on the nameserver.
func Starter(srv string, qtype string, hosts ...[]string) bool {
	tmp := []string{}
	for i := 0; i < len(hosts); i++ {
		tmp = append(tmp, hosts[i]...)
	}
	// if the slice has zero entries return false, there is no configured host to preload.
	if len(tmp) == 0 {
		return false
	}
	fmt.Printf("Preloading Nameserver: %s with query type: %s for domains: %s\n", srv, qtype, strings.Join(tmp, ", "))
	return true
}

// Hosts preload the nameserver with IP addresses for a given list of hostnames.
func (p *Preload) Hosts(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		start := time.Now()
		result, err := p.resolver.LookupIPAddr(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], aType, time.Now().Sub(start), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// CNAME preload the nameserver with CNAME lookups for a given list of hostnames.
func (p *Preload) CNAME(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		start := time.Now()
		result, err := p.resolver.LookupCNAME(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], cnameType, time.Now().Sub(start), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// NS preloads the nameserver records for a given list of hostnames.
func (p *Preload) NS(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		start := time.Now()
		result, err := p.resolver.LookupNS(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], nsType, time.Now().Sub(start), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// MX preloads the nameserver with the MX records for a given list of hostnames.
func (p *Preload) MX(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		start := time.Now()
		result, err := p.resolver.LookupMX(ctx, hosts[i])
		if err != nil {
			return err
		}

		err = p.Printer(hosts[i], mxType, time.Now().Sub(start), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// TXT preloads the nameserver with the TXT records for a given list of hostnames.
func (p *Preload) TXT(ctx context.Context, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		start := time.Now()
		result, err := p.resolver.LookupTXT(ctx, hosts[i])
		if err != nil {
			return err
		}
		err = p.Printer(hosts[i], txtType, time.Now().Sub(start), result)
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
				p.Hosts(context.Background(), str)
			}
		}
	case []*net.NS:
		for _, ns := range results.([]*net.NS) {
			str = append(str, ns.Host)
			if p.Full {
				p.Hosts(context.Background(), str)
			}
		}
	case []net.IPAddr:
		for _, ip := range results.([]net.IPAddr) {
			str = append(str, ip.IP.String())
		}
	default:
		return fmt.Errorf("ERROR: got type %+v\n", r.(string))
	}
	if !p.Quiet {
		fmt.Printf("Preloaded %s type %s in %s to %+s\n", hostname, qtype, duration, strings.Join(str, ", "))
	}
	return nil
}

func main() {
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
	fmt.Printf("Preload completed in %s\n", time.Now().Sub(start))
}
