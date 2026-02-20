package main

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	defaultJournalNumEntries = 20
	sdStateActive            = "active"
	sdStateActivating        = "activating"
	sdStateInactive          = "inactive"
	sdStateDeactivating      = "deactivating"
	sdStateFailed            = "failed"
	sdStateReloading         = "reloading"
	sdStateUnloaded          = "unloaded"
)

var (
	sdStatesAll = []string{sdStateActive, sdStateActivating, sdStateInactive, sdStateDeactivating, sdStateFailed, sdStateReloading}
	configLock  = new(sync.Mutex)
)

// Config contains service configuration
type Config struct {
	Units             []Unit `yaml:"units"`
	HTTP              HTTP   `yaml:"http"`
	SubscribeInterval int    `yaml:"subscribe_interval"`
	LogTimestamp      bool   `yaml:"log_timestamp"`
	JournalNumEntries uint64 `yaml:"journal_num_entries"`
}

func (c *Config) getUnit(name string) (Unit, error) {
	for _, unit := range c.Units {
		if unit.Name == name {
			return unit, nil
		}
	}

	return Unit{}, fmt.Errorf("unit %v not found in config", name)
}

func (c *Config) listUnits() []string {
	var units []string
	for _, unit := range c.Units {
		units = append(units, unit.Name)
	}

	return units
}

// Unit contains unit configuration
type Unit struct {
	Name      string   `yaml:"name"`
	OnActive  []string `yaml:"on_active"`
	OnInctive []string `yaml:"on_inactive"`
	OnFailed  []string `yaml:"on_failed"`
	BlockedBy []string `yaml:"blocked_by"`
}

// HTTP contains http service configuration
type HTTP struct {
	Enabled bool   `yaml:"enabled"`
	Bind    string `yaml:"bind"`
	XToken  string `yaml:"x_token"`
}

func loadConfig(path string) (*Config, error) {
	c := new(Config)
	err := updateConfig(c, path)
	if err != nil {
		return nil, err
	}

	// use default value if journal_num_entries is not set
	if c.JournalNumEntries == 0 {
		c.JournalNumEntries = defaultJournalNumEntries
	}

	return c, err
}

func updateConfig(cfg *Config, path string) error {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	configLock.Lock()
	defer configLock.Unlock()
	*cfg = c
	return nil
}
