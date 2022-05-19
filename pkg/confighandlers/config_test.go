package confighandlers

import (
	"testing"
)

func TestConfigurationPrintEmptyConfigration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Generate Empty Configuration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Configuration{}
			cfg.PrintEmptyConfigration()
		})
	}
}
