package main

import (
	"encoding/json"
	"log"
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
	defer func() {
		if err := file.Close(); err != nil {
			log.Panic(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&conf); err != nil {
		return nil, err
	}

	return conf, err
}
