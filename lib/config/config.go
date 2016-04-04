package config

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type AppConfig struct {
	// LibraryPath is the directory where the library is persisted
	LibraryPath string

	// TemplatePath is the directory containing the html templates and
	// static assets
	TemplatePath string

	// networkAddr the address the library webservice will listen on
	// e.g. ":8080"
	NetworkAddr string
}

func LoadConfigFromFile(configFile string) (*AppConfig, error) {
	errMsg := "error loading config file %v"
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	config := &AppConfig{}
	err = json.Unmarshal(configData, config)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	err = validateConfig(config)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return config, nil
}

func validateConfig(config *AppConfig) error {
	if config.TemplatePath == "" {
		return fmt.Errorf("Missing template path")
	}
	if config.LibraryPath == "" {
		return fmt.Errorf("Missing library path")
	}
	if config.NetworkAddr == "" {
		return fmt.Errorf("Missing network address")
	}
	return nil
}