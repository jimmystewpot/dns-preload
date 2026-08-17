package confighandlers

import (
	"fmt"
	"strings"
	"testing"

	validator "github.com/go-playground/validator/v10"
)

func TestLoadConfigFromFile(t *testing.T) {
	type args struct {
		cfgfile *string
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryType
		wantErr bool
	}{
		{
			name: "test configuration file",
			args: args{
				cfgfile: ptr("test_data/complete_config_sample.yaml"),
			},
			want: &QueryType{
				Hosts: []string{
					"google.com",
					"microsoft.com",
				},
				NS: []string{
					"github.com",
					"bitbucket.org",
				},
				MX: []string{
					"gmail.com",
					"hotmail.com",
				},
				TXT: []string{
					"salesforce.com",
					"linkedin.com",
				},
			},
			wantErr: false,
		},
		{
			name: "missing test configuration file",
			args: args{
				cfgfile: ptr("test_data/missing_file_example.yaml"),
			},
			want:    &QueryType{},
			wantErr: true,
		},
		{
			name: "bad test configuration file",
			args: args{
				cfgfile: ptr("test_data/bad_configuration_sample.yaml"),
			},
			want:    &QueryType{},
			wantErr: true,
		},
		{
			name: "too large configuration file",
			args: args{
				cfgfile: ptr("test_data/overrun_configuration_example.yaml"),
			},
			want:    &QueryType{},
			wantErr: true,
		},
		{
			name: "missing configuration keys",
			args: args{
				cfgfile: ptr("test_data/missing_config_keys.yaml"),
			},
			want:    &QueryType{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfigFromFile(tt.args.cfgfile)
			fmt.Printf("%+v\n", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.QueryType.Hosts[0] != tt.want.Hosts[0] {
				t.Errorf("LoadConfigFromFile() expected %s got %s", tt.want.Hosts[0], got.QueryType.Hosts[0])
			}
		})
	}
}

func TestQueryListPopulateCounts(t *testing.T) {
	type fields struct {
		QueryType QueryType
	}
	tests := []struct {
		name    string
		fields  fields
		want    fields
		wantErr bool
	}{
		{
			name: "Simple test",
			fields: fields{
				QueryType{
					Cname: []string{"one", "two"},
				},
			},
			want: fields{
				QueryType{
					CnameCount: uint16(2),
				},
			},
		},
		{
			name: "zero test",
			fields: fields{
				QueryType{
					Cname: []string{},
				},
			},
			want: fields{
				QueryType{
					CnameCount: uint16(0),
				},
			},
		},
		{
			name: "overload test cname",
			fields: fields{
				QueryType{
					Cname: garbage(100000000),
				},
			},
			want: fields{
				QueryType{
					CnameCount: uint16(0),
				},
			},
			wantErr: true,
		},
		{
			name: "overload test hosts",
			fields: fields{
				QueryType{
					Hosts: garbage(100000000),
				},
			},
			want: fields{
				QueryType{
					HostsCount: uint16(0),
				},
			},
			wantErr: true,
		},
		{
			name: "overload test mx",
			fields: fields{
				QueryType{
					MX: garbage(100000000),
				},
			},
			want: fields{
				QueryType{
					MXCount: uint16(0),
				},
			},
			wantErr: true,
		},
		{
			name: "overload test ns",
			fields: fields{
				QueryType{
					NS: garbage(100000000),
				},
			},
			want: fields{
				QueryType{
					NSCount: uint16(0),
				},
			},
			wantErr: true,
		},
		{
			name: "overload test txt",
			fields: fields{
				QueryType{
					TXT: garbage(100000000),
				},
			},
			want: fields{
				QueryType{
					TXTCount: uint16(0),
				},
			},
			wantErr: true,
		},
		{
			name: "overload test ptr",
			fields: fields{
				QueryType{
					PTR: garbage(100000000),
				},
			},
			want: fields{
				QueryType{
					PTRCount: uint16(0),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ql := &Configuration{
				QueryType: tt.fields.QueryType,
			}
			err := ql.PopulateCounts()
			if (err != nil) && (tt.wantErr == false) {
				t.Errorf("got error %s when we did not expect it", err)
			}
			if ql.QueryType.CnameCount != tt.want.QueryType.CnameCount {
				t.Errorf("expected %d got %d", tt.want.QueryType.CnameCount, ql.QueryType.CnameCount)
			}
		})
	}
}

func TestCount(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				s: garbage(2),
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "error path",
			args: args{
				s: garbage(100000),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := count(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func garbage(s int64) []string {
	x := make([]string, s)
	return x
}

func ptr(s string) *string {
	return &s
}

func TestUint16Boundaries(t *testing.T) {
	t.Run("boundary 65535", func(t *testing.T) {
		s := make([]string, 65535)
		got, err := Uint16(s)
		if err != nil {
			t.Errorf("Uint16(65535) unexpected error = %v", err)
		}
		if got != 65535 {
			t.Errorf("Uint16(65535) = %d, want 65535", got)
		}
	})

	t.Run("boundary 65536 overflow", func(t *testing.T) {
		s := make([]string, 65536)
		_, err := Uint16(s)
		if err == nil {
			t.Error("Uint16(65536) expected error for uint16 overflow, got nil")
		}
	})
}

func TestLoadConfigValidationFailures(t *testing.T) {
	t.Run("invalid fqdn in hosts", func(t *testing.T) {
		yamlData := `
query_type:
  hosts:
    - "invalid fqdn with spaces"
`
		var cfg Configuration
		err := cfg.LoadConfig(strings.NewReader(yamlData))
		if err != nil {
			t.Fatalf("LoadConfig error = %v", err)
		}
		_ = cfg.PopulateCounts()
		val := validator.New()
		err = val.Struct(&cfg)
		if err == nil {
			t.Error("expected validation error for invalid FQDN, got nil")
		}
	})

	t.Run("invalid ip in ptr", func(t *testing.T) {
		yamlData := `
query_type:
  ptr:
    - "999.999.999.999"
`
		var cfg Configuration
		err := cfg.LoadConfig(strings.NewReader(yamlData))
		if err != nil {
			t.Fatalf("LoadConfig error = %v", err)
		}
		_ = cfg.PopulateCounts()
		val := validator.New()
		err = val.Struct(&cfg)
		if err == nil {
			t.Error("expected validation error for invalid IP in PTR, got nil")
		}
	})
}
