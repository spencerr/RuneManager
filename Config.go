package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	DatabaseConnectionString	string 	`json:"connection_string"`
	Environment 				string	`json:"env"`
}

var (
	config *Config
)

func loadConfig() {
	
	if data, err := ioutil.ReadFile("config.json"); err == nil {
		if err := json.Unmarshal(data, &config); err == nil {
			fmt.Printf("Successfully loaded configurations. Environment set to %s", config.Environment)
			return
		}
	}

	panic(err)	
}