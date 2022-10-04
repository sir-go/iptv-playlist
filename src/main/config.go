package main

import (
	"encoding/json"
	"os"
)

type Cfg struct {
	Service struct {
		Port int `json:"port"`
	} `json:"service"`

	Database struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"database"`

	UnicastUrl string `json:"unicast_url"`
}

func LoadConfig(confpath string) (*Cfg, error) {

	conf := new(Cfg)
	file, err := os.Open(confpath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(&conf)

	return conf, err
}
