package confighandlers

import (
	"testing"
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
			name: "bad test configuration file",
			args: args{
				cfgfile: ptr("test_data/bad_config_sample.yaml"),
			},
			want:    &QueryType{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfigFromFile(tt.args.cfgfile)
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

func ptr(s string) *string {
	return &s
}
