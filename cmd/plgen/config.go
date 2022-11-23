package main

import (
	"github.com/go-yaml/yaml"
)

type (
	cfgService struct {
		Port int `yaml:"port"`
	}

	cfgDb struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}

	cfgTemplates struct {
		M3U  string `yaml:"m3u"`
		XSPF string `yaml:"xspf"`
	}

	cfgPlaylist struct {
		Templates  cfgTemplates `yaml:"templates"`
		UnicastUrl string       `yaml:"unicast_url"`
	}

	Config struct {
		Service  cfgService  `yaml:"service"`
		Db       cfgDb       `yaml:"db"`
		Playlist cfgPlaylist `yaml:"playlist"`
	}
)

func LoadConfig(b []byte) (cfg *Config, err error) {
	cfg = new(Config)
	if err = yaml.UnmarshalStrict(b, cfg); err != nil {
		return nil, err
	}
	return
}
