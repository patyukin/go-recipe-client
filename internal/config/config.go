package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port       int    `yaml:"port" validate:"required"`
	RemotePort int    `yaml:"remote_port" validate:"required"`
	RemoteHost string `yaml:"remote_host" validate:"required"`
}

func LoadConfig() (*Config, error) {
	yamlConfigFilePath := os.Getenv("YAML_CONFIG_FILE_PATH")
	if yamlConfigFilePath == "" {
		return nil, fmt.Errorf("yaml config file path is not set")
	}

	f, err := os.Open(yamlConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Printf("unable to close config file: %v", err)
		}
	}(f)

	cfg := &Config{}

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	validate := validator.New()
	if err = validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
