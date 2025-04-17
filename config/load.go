package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	lock           sync.Mutex
	isSet          bool
	singleInstance Config
	configPath     string
)

func UpdatePath(path string) {
	configPath = path
}

func Get() Config {
	lock.Lock()
	defer lock.Unlock()
	if !isSet {
		config, err := initializeConfig(configPath)
		if err != nil {
			panic(err)
		}
		singleInstance = config
    isSet = true
	}
	return singleInstance
}

func update(config Config) {
	lock.Lock()
	defer lock.Unlock()
	singleInstance = config
}

func initializeConfig(path string) (Config, error) {
	var config Config
	err := ParseYamlFile(path, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}

func ParseYamlFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}
