package confighandlers

import (
	"testing"
)

func TestConfigurationPrintEmptyConfigration(t *testing.T) {
	tests := []struct {
		name  string
		quiet bool
	}{
		{
			name:  "Generate Empty Configuration",
			quiet: false,
		},
		{
			name:  "Generate Empty Configuration quiet",
			quiet: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Configuration{}
			err := cfg.PrintEmptyConfigration(tt.quiet)
			if err != nil {
				t.Errorf("error in PrintEmptyconfiguration %s", err)
			}
		})
	}
}
