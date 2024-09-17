package confighandlers

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

const (
	Mx         string = "mx"
	Ns         string = "ns"
	Txt        string = "txt"
	Cname      string = "cname"
	Hosts      string = "hosts"
	Ptr        string = "ptr"
	nilRecords uint16 = 0
)

var (
	// queryTypes is used to iterate through all of the commands when the all cmd is used.
	QueryTypes = []string{Hosts, Cname, Mx, Ns, Txt, Ptr}
)

type Configuration struct {
	QueryType QueryType `yaml:"query_type" json:"query_type" validate:"required"`
}

// QueryType lsits out the structure for the different domains and their query type
type QueryType struct {
	// Cnames for doing a cname lookup
	Cname []string `yaml:"cname" json:"cname" validate:"dive,fqdn"`
	// Hosts for doing a query for type A and AAAA
	Hosts []string `yaml:"hosts" json:"hosts" validate:"dive,fqdn"`
	// ns for doing a query for type NS
	NS []string `yaml:"ns" json:"ns" validate:"dive,fqdn"`
	// MX for doing a query for type MX
	MX []string `yaml:"mx" json:"mx" validate:"dive,fqdn"`
	// TXT for doing a query for type TXT
	TXT []string `yaml:"txt" json:"txt" validate:"dive,fqdn"`
	// PTR for doing a query for type PTR
	PTR []string `yaml:"ptr" json:"ptr" validate:"dive,ip_addr"`
	// Metrics values below this point.
	CnameCount uint16 `yaml:",omitempty"`
	HostsCount uint16 `yaml:",omitempty"`
	NSCount    uint16 `yaml:",omitempty"`
	MXCount    uint16 `yaml:",omitempty"`
	TXTCount   uint16 `yaml:",omitempty"`
	PTRCount   uint16 `yaml:",omitempty"`
}

// PopulateCounts for how many domains are in each query_type.
func (cfg *Configuration) PopulateCounts() error {
	var err error
	cfg.QueryType.CnameCount, err = count(cfg.QueryType.Cname)
	if err != nil {
		return err
	}
	cfg.QueryType.HostsCount, err = count(cfg.QueryType.Hosts)
	if err != nil {
		return err
	}
	cfg.QueryType.NSCount, err = count(cfg.QueryType.NS)
	if err != nil {
		return err
	}
	cfg.QueryType.MXCount, err = count(cfg.QueryType.MX)
	if err != nil {
		return err
	}
	cfg.QueryType.TXTCount, err = count(cfg.QueryType.TXT)
	if err != nil {
		return err
	}
	cfg.QueryType.PTRCount, err = count(cfg.QueryType.PTR)
	if err != nil {
		return err
	}
	return nil
}

// count wrapper for uint16
func count(s []string) (uint16, error) {
	return Uint16(s)
}

// Uint16 safely check that the conversion to uint16 works.
func Uint16(s []string) (uint16, error) {
	size := len(s)
	if size > 65535 || size < 0 {
		return 0, fmt.Errorf("%d not convertable to uint16", size)
	}
	return uint16(size), nil
}

// PrintEmptyConfiguration is used to generate an empty configuration to stdout
func (cfg *Configuration) PrintEmptyConfigration(quiet bool) error {
	yamlConfiguration, err := yaml.Marshal(&Configuration{})
	if err != nil {
		return err
	}
	if !quiet {
		fmt.Printf("---\n%s\n\n", string(yamlConfiguration))
	}
	return nil
}
