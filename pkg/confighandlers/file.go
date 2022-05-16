package confighandlers

import (
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	Mx    string = "mx"
	Ns    string = "ns"
	Txt   string = "txt"
	Cname string = "cname"
	Hosts string = "hosts"
)

var (
	// queryTypes is used to iterate through all of the commands when the all cmd is used.
	QueryTypes = []string{Hosts, Cname, Mx, Ns, Txt}
)

type Configuration struct {
	QueryType QueryType `yaml:"query_type" json:"query_type"`
}

// QueryType lsits out the structure for the different domains and their query type
type QueryType struct {
	// Cnames for doing a cname lookup
	Cname      []string `yaml:"cname" json:"cname"`
	CnameCount uint16
	// Hosts for doing a query for type A and AAAA
	Hosts      []string `yaml:"hosts" json:"hosts"`
	HostsCount uint16
	// ns for doing a query for type NS
	NS      []string `yaml:"ns" json:"ns"`
	NSCount uint16
	// MX for doing a query for type MX
	MX      []string `yaml:"mx" json:"mx"`
	MXCount uint16
	// TXT for doing a query for type TXT
	TXT      []string `yaml:"txt" json:"txt"`
	TXTCount uint16
}

// ReadConfig wil read the YAML file from disk and render it into the DomainConfig struct.
func (cfg *Configuration) LoadConfig(r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

// PopulateCounts for how many domains are in each query_type.
func (cfg *Configuration) PopulateCounts() error {
	var err error
	cfg.QueryType.CnameCount, err = count(cfg.QueryType.Cname)
	fmt.Println(cfg.QueryType.CnameCount)
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
	return nil
}

// count wrapper for uint16
func count(s []string) (uint16, error) {
	return Uint16(s)
}

// loadConfig will load the configuration from file.
func LoadConfigFromFile(cfgfile *string) (*Configuration, error) {
	// cfg is a slice of strings unmarsalled from YAML
	cfg := new(Configuration)

	// Open passed in filename
	f, err := os.Open(*cfgfile)
	if err != nil {
		return &Configuration{}, err
	}
	defer f.Close()

	// Load from reader, validate
	err = cfg.LoadConfig(f)
	if err != nil {
		return &Configuration{}, err
	}
	err = cfg.PopulateCounts()
	if err != nil {
		return &Configuration{}, err
	}

	return cfg, nil
}

// Uint16 safely check that the conversion to uint16 works.
func Uint16(s []string) (uint16, error) {
	size := len(s)
	if size > 65535 || size < 0 {
		return 0, fmt.Errorf("%d not convertable to uint16", size)
	}
	return uint16(size), nil
}
