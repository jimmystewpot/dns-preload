package confighandlers

import (
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	// queryTypes is used to iterate through all of the commands when the all cmd is used.
	QueryTypes = []string{"hosts", "cname", "mx", "ns", "txt"}
)

type QueryList struct {
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
func (ql *QueryList) LoadConfig(r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(ql)
	if err != nil {
		return err
	}
	return nil
}

func (ql *QueryList) PopulateCounts() error {
	var err error
	ql.QueryType.CnameCount, err = count(ql.QueryType.Cname)
	if err != nil {
		return err
	}
	ql.QueryType.HostsCount, err = count(ql.QueryType.Hosts)
	if err != nil {
		return err
	}
	ql.QueryType.NSCount, err = count(ql.QueryType.NS)
	if err != nil {
		return err
	}
	ql.QueryType.MXCount, err = count(ql.QueryType.MX)
	if err != nil {
		return err
	}
	ql.QueryType.TXTCount, err = count(ql.QueryType.TXT)
	if err != nil {
		return err
	}
	return nil
}

func count(s []string) (uint16, error) {
	return Uint16(s)
}

// loadConfig will load the configuration from file.
func LoadConfigFromFile(cfgfile *string) (*QueryList, error) {
	// cfg is a slice of strings unmarsalled from YAML
	cfg := new(QueryList)

	// Open passed in filename
	f, err := os.Open(*cfgfile)
	if err != nil {
		return &QueryList{}, err
	}
	defer f.Close()

	// Load from reader, validate
	err = cfg.LoadConfig(f)
	if err != nil {
		return &QueryList{}, err
	}
	err = cfg.PopulateCounts()
	if err != nil {
		return &QueryList{}, err
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
