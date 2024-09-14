package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var ConfigFile string

func New() *Config {
	var c Config

	file, err := os.Open(ConfigFile)
	if err != nil {
		return &c
	}

	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		return &c
	}

	return &c
}
