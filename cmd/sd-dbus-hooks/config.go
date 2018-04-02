package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Units            []Unit `yaml:"units"`
	HTTP             HTTP   `yaml:"http"`
	SubscribeTimeout int    `yaml:"subscribe_timeout"`
}

func (c *Config) getUnit(name string) (Unit, error) {
	for _, unit := range c.Units {
		if unit.Name == name {
			return unit, nil
		}
	}

	return Unit{}, fmt.Errorf("unit %v not found in config", name)
}

type Unit struct {
	Name      string   `yaml:"name"`
	OnActive  []string `yaml:"on_active"`
	OnInctive []string `yaml:"on_inactive"`
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
