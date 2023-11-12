package confighandlers

import (
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	yaml "gopkg.in/yaml.v3"
)

var (
	validate *validator.Validate
)

// ReadConfig wil read the YAML file from disk and render it into the DomainConfig struct.
func (cfg *Configuration) LoadConfig(r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

// loadConfig will load the configuration from file.
//
//nolint:lll // validation has many fields that can't be zero.
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
	validate = validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return &Configuration{}, err
	}
	if (cfg.QueryType.CnameCount == 0) && (cfg.QueryType.HostsCount == 0) && (cfg.QueryType.MXCount == 0) && (cfg.QueryType.PTRCount == 00) && (cfg.QueryType.TXTCount == 0) {
		return &Configuration{}, fmt.Errorf("empty configuration or invalid keys")
	}

	return cfg, nil
}
