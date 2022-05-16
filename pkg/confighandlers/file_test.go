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
