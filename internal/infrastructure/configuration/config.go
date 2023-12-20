package configuration

import (
	"ara-server/util/log"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const keyENV = "ENV"

func (config *Config) GetConfig() AppConfig {
	return config.config
}

func InitializeConfig() (Config, error) {
	env := os.Getenv(keyENV)
	if env == "" {
		env = "development"
	}
	log.Info(fmt.Sprintf("Running on Environment %s", env))

	workingdir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	configPath := fmt.Sprintf("%s/config/config.%s.yaml", workingdir, env)
	yamlFile, err := ioutil.ReadFile(configPath)
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
	}, nil
}
