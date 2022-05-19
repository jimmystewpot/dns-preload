package confighandlers

import (
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
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
