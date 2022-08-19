package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DBConfig     `yaml:"database-settings"`
	Server   ServerConfig `yaml:"server-settings"`
}

type DBConfig struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
}

// Load decodes a yaml configuration file
func Load(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Config{}
	d := yaml.NewDecoder(f)

	if err := d.Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
