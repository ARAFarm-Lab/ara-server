package configuration

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const keyENV = "ENV"

func (config *Config) GetConfig() AppConfig {
	return config.config
}

func (config *Config) IsDevelopment() bool {
	return config.env == Development
}

func InitializeConfig() (Config, error) {
	env := strings.ToLower(os.Getenv(keyENV))
	if env == "" {
		env = "development"
	}

	workingdir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	configPath := fmt.Sprintf("%s/config/config.%s.yaml", workingdir, env)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config AppConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}

	return Config{
		config: config,
		env:    ENV(env),
	}, nil
}
