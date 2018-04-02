package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Units []Unit `yaml:"units"`
	HTTP  HTTP   `yaml:"http"`
}

type Unit struct {
	OnActive  []string `yaml:"on_active"`
	OnFailed  []string `yaml:"on_failed"`
	BlockedBy []string `yaml:"blocked_by"`
}

type HTTP struct {
	Bind string `yaml:"bind"`
}

func LoadConfig(path string) (*Config, error) {
	var c Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, err
}
