package config

import (
	"fmt"
	"qbittorrent_exporter/lib/log"
	"qbittorrent_exporter/lib/parser"
	"qbittorrent_exporter/validator"
	"strconv"
	"sync"
)

var (
	configPath     string = "config.yaml"
	isSet          bool   = false
	singleInstance Config
	lock           sync.Mutex
)

type (
	Config struct {
		QBittorrent QBittorrentConfig `yaml:"qBittorrent"`
		Metrics     MetricsConfig     `yaml:"metrics"`
		Global      GlobalConfig      `yaml:"global"`
	}
	QBittorrentConfig struct {
		BaseURL  string `yaml:"baseUrl" env:"QBE_URL"`
		Username string `yaml:"username" env:"QBE_USERNAME"`
		Password string `yaml:"password" env:"QBE_PASSWORD"`
	}
	MetricsConfig struct {
		Port    string `yaml:"port" env:"QBE_METRICS_PORT"`
		UrlPath string `yaml:"urlPath" env:"QBE_METRICS_PATH"`
	}
	GlobalConfig struct {
		StatePath string `yaml:"statePath" env:"QBE_STATE_PATH"`
	}
)

func UpdatePath(path string) {
	configPath = path
}

func Get() Config {
	lock.Lock()
	defer lock.Unlock()
	if !isSet {
		log.Debug("Config isSet=false, initializing...")
		singleInstance = initializeConfig(configPath)
		isSet = true
	}
	return singleInstance
}

func initializeConfig(path string) Config {
	var config Config
	if err := validator.ValidatePath(configPath, false); err != nil {
		panic(err)
	}
	if err := parser.ParseYamlFile(path, &config); err != nil {
		panic(fmt.Errorf("error parsing config: %w", err))
	}
	log.Debug("Loading environment variables into the config")
	loadEnvs(&config)
	return config
}

// TODO:
func Validate(config Config) error {
	if _, err := strconv.ParseUint(config.Metrics.Port, 10, 16); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to validate config.Metrics.Port")
	}
	return nil
}
