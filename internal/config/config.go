package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type Config struct {
	HTTPServerListenAddr string `yaml:"httpServerListenAddr"`
	DBSourceName         string `yaml:"dbSourceName"`
}

var (
	cfg    *Config
	cfgErr error
	once   sync.Once
)

func New(configPath string) (*Config, error) {
	once.Do(func() {
		cfg, cfgErr = parseConfig(configPath)
	})

	return cfg, cfgErr
}

func parseConfig(configPath string) (*Config, error) {
	var c Config

	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(rawConfig, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
