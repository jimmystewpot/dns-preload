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
	"github.com/jimmystewpot/dns-preload/pkg/dns"
)

const (
	// these const strings are used to store the DNS query types for reuse.
	queryTypeAStr     string = "A, AAAA"
	queryTypeCNAMEStr string = "CNAME"
	queryTypeMXStr    string = "MX"
	queryTypeNSStr    string = "NS"
	queryTypeTXTStr   string = "TXT"
	// print messages that are used more than once.
	infoMessage          string = "Preloading Nameserver: %s with query type: %s for domains: %s\n"
	qTypeEmptyErrMessage string = "Preloading error: query type %s has no entries in the configuration"
	qTypeErrMessage      string = "preloading error: query type %s is not a valid query type"
	completedMessage     string = "Preload completed in %s"
)

var (
	cli struct {
		All    Preload       `cmd:"" help:"preload all of the following types from the configuration file"`
		Cname  Preload       `cmd:"" help:"preload only the cname entries from the configuration file"`
		Hosts  Preload       `cmd:"" help:"preload only the hosts entries from the configuration file, this does an A and AAAA lookup"`
		Mx     Preload       `cmd:"" help:"preload only the mx entries from the configuration file"`
		Ns     Preload       `cmd:"" help:"preload only the ns entries from the configuration file"`
		Txt    Preload       `cmd:"" help:"preload only the txt entries from the configuration file"`
		Config Config        `cmd:"" help:"generate an empty configuration file to stdout"`
		Delay  time.Duration `default:"0s" help:"How long to wait until the queries are executed"`
	}
	start time.Time
	quiet bool
)

type Preload struct {
	ConfigFile string        `required:"" help:"The configuration file to read the domain list to query from"`
	Server     string        `default:"localhost" help:"The server to query to seed the domain list into"`
	Port       string        `default:"53" help:"The port the DNS server listens for requests on"`
	Workers    uint8         `default:"1" help:"The number of concurrent goroutines used to query the DNS server (not implemented)"`
	Mute       bool          `default:"false" help:"Suppress the preload task output to the console"`
	Quiet      bool          `default:"false" help:"Suppress the preload response output to the console"`
	Full       bool          `default:"true" help:"For record types that return a Hostname ensure that these are resolved"`
	Debug      bool          `default:"false" help:"Debug mode"`
	Timeout    time.Duration `default:"30s" help:"The timeout for DNS queries to succeed"`
	resolver   dns.Resolver
	nameserver string
}

type Config struct {
	Quiet    bool `default:"false" help:"Suppress the info output to the console"`
	Generate struct {
		Generate bool `default:"true" help:"Generate an empty configuration and output it to stdout"`
	} `cmd:"" help:"Generate a configuration file"`
	Validate struct {
		ConfigFile string `required:"" help:"The configuration file to load"`
	} `cmd:"" help:"Validate a configuration file"`
}

// Config Run() prints an empty YAML configuration to stdout.
func (c *Config) Run(cmd string) error {
	switch cmd {
	case "config generate":
		cfg := confighandlers.Configuration{}
		return cfg.PrintEmptyConfigration(c.Quiet)
	case "config validate":
		return fmt.Errorf("not yet implemented")
	}
	return fmt.Errorf("unknown command %s", cmd)
}

func (p *Preload) Run(cmd string) error {
	quiet = p.Quiet
	cfg, err := confighandlers.LoadConfigFromFile(&p.ConfigFile)
	if err != nil {
		return err
	}

	p.nameserver = net.JoinHostPort(p.Server, p.Port)
	p.resolver = dns.NewResolver(p.nameserver, p.Timeout)

	ctx := context.Background()
	return p.RunQueries(ctx, cmd, cfg)
}

// RunQueries breaks out the command switch statement allowing me to write better tests by adding a mock resolver.
func (p *Preload) RunQueries(ctx context.Context, cmd string, cfg *confighandlers.Configuration) error {
	switch cmd {
	case confighandlers.Cname:
		if cfg.QueryType.CnameCount != 0 {
			p.IntroPrinter(queryTypeCNAMEStr, cfg.QueryType.Cname)
			return p.CNAME(ctx, cfg.QueryType.Cname)
		}
	case confighandlers.Hosts:
		if cfg.QueryType.HostsCount != 0 {
			p.IntroPrinter(queryTypeAStr, cfg.QueryType.Hosts)
			return p.Hosts(ctx, cfg.QueryType.Hosts)
		}
	case confighandlers.Mx:
		if cfg.QueryType.MXCount != 0 {
			p.IntroPrinter(queryTypeMXStr, cfg.QueryType.MX)
			return p.MX(ctx, cfg.QueryType.MX)
		}
	case confighandlers.Ns:
		if cfg.QueryType.NSCount != 0 {
			p.IntroPrinter(queryTypeNSStr, cfg.QueryType.NS)
			return p.NS(ctx, cfg.QueryType.NS)
		}
	case confighandlers.Txt:
		if cfg.QueryType.TXTCount != 0 {
			p.IntroPrinter(queryTypeTXTStr, cfg.QueryType.TXT)
			return p.TXT(ctx, cfg.QueryType.TXT)
		}
	default: // no known query type fallback error handling.
		return fmt.Errorf(qTypeErrMessage, cmd)
	}
	if p.Debug {
		fmt.Printf(qTypeEmptyErrMessage+"\n", cmd)
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
		err = p.ResultsPrinter(hosts[i], queryTypeAStr, time.Since(s), result)
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
		err = p.ResultsPrinter(hosts[i], queryTypeCNAMEStr, time.Since(s), result)
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
		err = p.ResultsPrinter(hosts[i], queryTypeNSStr, time.Since(s), result)
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

		err = p.ResultsPrinter(hosts[i], queryTypeMXStr, time.Since(s), result)
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
		err = p.ResultsPrinter(hosts[i], queryTypeTXTStr, time.Since(s), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// String provides output to the console for the results of the preloading.
func (p *Preload) ResultsPrinter(hostname string, qtype string, duration time.Duration, results interface{}) error {
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
		}
	case []*net.NS:
		for _, ns := range results.([]*net.NS) {
			str = append(str, ns.Host)
		}
	case []net.IPAddr:
		for _, ip := range results.([]net.IPAddr) {
			str = append(str, ip.IP.String())
		}
	default:
		return fmt.Errorf("error: unknown type %T", r)
	}
	if p.Full {
		// mx and ns record types return hostnames, if full is on we should resolve the final targets.
		if (qtype == queryTypeMXStr) || (qtype == queryTypeNSStr) {
			err := p.Hosts(context.Background(), str)
			if err != nil {
				return err
			}
		}
	}
	if !p.Quiet {
		fmt.Printf("Preloaded %s type %s in %s to %+s\n", hostname, qtype, duration, strings.Join(str, ", "))
	}

	return nil
}

// IntroPrinter outputs the info on what domains and servers are being reloaded.
func (p *Preload) IntroPrinter(queryType string, hosts []string) {
	if !p.Mute {
		fmt.Printf(infoMessage+"\n", p.nameserver, queryType, strings.Join(hosts, ", "))
	}
}

// CompletedPrinter prints out a completion message and a timer to stdout.
func completedPrinter(quiet bool, t time.Time) {
	if !quiet {
		fmt.Printf(completedMessage+"\n", time.Since(start))
	}
}

func main() {
	start = time.Now()
	cmd := kong.Parse(&cli,
		kong.Name(os.Args[0]),
		kong.Description("Preload a DNS cache with a list of hostnames from a YAML configuration file."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	time.Sleep(cli.Delay)
	switch cmd.Command() {
	case "all":
		e := make([]error, 0)
		for _, queryType := range confighandlers.QueryTypes {
			err := cmd.Run(queryType)
			if err != nil {
				e = append(e, fmt.Errorf("quertyType: %s, err: %s", queryType, err))
			}
		}
		completedPrinter(quiet, start)

		if len(e) != 0 {
			var errLog error
			for _, errLog = range e {
				fmt.Printf("%s ", errLog)
			}
			cmd.FatalIfErrorf(errLog)
		}
	default:
		err := cmd.Run(cmd.Command())
		completedPrinter(quiet, start)
		cmd.FatalIfErrorf(err)
	}
}
