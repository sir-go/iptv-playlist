package main

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	const confData = `
service:
  port: 4444

db:
  host: db-host
  port: 3306
  user: db-user
  password: db-password
  dbname: db-name

playlist:
  templates:
    m3u: /path/to/templates/m3u.tmpl
    xspf: /path/to/templates/xspf.tmpl
  unicast_url: http://iptv-stream-url:8000
`

	confGood := Config{
		Service: cfgService{Port: 4444},
		Db: cfgDb{
			Host:     "db-host",
			Port:     3306,
			User:     "db-user",
			Password: "db-password",
			Dbname:   "db-name",
		},
		Playlist: cfgPlaylist{
			Templates: cfgTemplates{
				M3U:  "/path/to/templates/m3u.tmpl",
				XSPF: "/path/to/templates/xspf.tmpl",
			},
			UnicastUrl: "http://iptv-stream-url:8000",
		},
	}

	type args struct {
		b []byte
	}

	tests := []struct {
		name    string
		args    args
		wantCfg *Config
		wantErr bool
	}{
		{"empty", args{[]byte{}}, &Config{}, false},
		{"ok", args{[]byte(confData)}, &confGood, false},
		{"bad", args{[]byte(confData + "\nextra: field")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := LoadConfig(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("LoadConfig() gotCfg = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}
