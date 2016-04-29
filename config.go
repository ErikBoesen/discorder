package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Email     string `json:"email"`
	AuthToken string `json:"authToken"` // Not used currently, but planned

	AllPrivateMode    bool     `json:"allPrivateMode"`
	LastChannel       string   `json:"lastChannel"`
	ListeningChannels []string `json:"listeningChannels"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var decoded Config
	err = json.Unmarshal(file, &decoded)
	return &decoded, err
}

func (c *Config) Save(path string) error {
	eencoded, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, eencoded, os.FileMode(0755))
}
