package main

import (
	"log"
	"encoding/json"
	"io/ioutil"
	)

type WeebConfig struct {
	Name, Host, Env string
	Port      uint32
	Registration *WeebRegistrationConfig
}

type WeebRegistrationConfig struct {
	Enabled     bool
	Host, Token string
}

func LoadConfig() WeebConfig {
	configFile, err := ioutil.ReadFile("config/main.json")
	if err != nil {
		log.Fatal(err)
	}
	config := WeebConfig{}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
