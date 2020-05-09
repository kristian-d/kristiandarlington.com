package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Original string
	Port string
	ReadTimeout uint16 `yaml:"readTimeout"`
	WriteTimeout uint16 `yaml:"writeTimeout"`
	MaxConnections uint16 `yaml:"maxConnections"`
}

var (
	DefaultConfig = Config{
		Original: "",
		Port: "8080",
		ReadTimeout: 5000,
		WriteTimeout: 5000,
		MaxConnections: 128,
	}
)

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	// if the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.
	*cfg = DefaultConfig

	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.Original = s
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfgAsString, err := Load(string(content))
	if err != nil {
		return nil, err
	}
	return cfgAsString, nil
}
