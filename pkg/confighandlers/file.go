package confighandlers

import (
	"io"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type QueryList struct {
	QueryType QueryType `yaml:"query_type" json:"query_type"`
}

// QueryType lsits out the structure for the different domains and their query type
type QueryType struct {
	// Cnames for doing a cname lookup
	Cname []string `yaml:"cname" json:"cname"`
	// Hosts for doing a query for type A and AAAA
	Hosts []string `yaml:"hosts" json:"hosts"`
	// ns for doing a query for type NS
	NS []string `yaml:"ns" json:"ns"`
	// MX for doing a query for type MX
	MX []string `yaml:"mx" json:"mx"`
	// TXT for doing a query for type TXT
	TXT []string `yaml:"txt" json:"txt"`
}

// ReadConfig wil read the YAML file from disk and render it into the DomainConfig struct.
func (qt *QueryList) LoadConfig(r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(qt)
	if err != nil {
		return err
	}
	return nil
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

	return cfg, nil
}
