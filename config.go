package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database    DBConfig     `yaml:"database"`
	LogFile     string       `yaml:"log-file"`
	Server      ServerConfig `yaml:"server-settings"`
	SwaggerSpec string       `yaml:"swagger-spec"`
	// TODO: change Prod field to deployment mode field
	Prod bool
}

type DBConfig struct {
	Path   string `yaml:"path"`
	IsTest bool   `yaml:"is-test"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
}

// loadConfig decodes a yaml configuration file, and sets deployment mode
func loadConfig(prodMode bool, configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// decode server config
	conf := &Config{}
	d := yaml.NewDecoder(f)

	if err := d.Decode(conf); err != nil {
		return nil, err
	}

	conf.Prod = prodMode

	return conf, nil
}
