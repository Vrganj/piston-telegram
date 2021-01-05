package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config - config.json
type Config struct {
	Token string `json:"token"`
}

func loadConfig() (*Config, error) {
	text, err := ioutil.ReadFile("config.json")

	if err != nil {
		return nil, err
	}

	config := Config{}

	if err := json.Unmarshal(text, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
