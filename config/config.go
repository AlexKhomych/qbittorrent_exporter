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
	configPath string = "config.yaml"
	isLoaded   bool   = false
	instance   Config
	mu         sync.Mutex
)

type Config struct {
	QBittorrent QBittorrentConfig `yaml:"qBittorrent"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	Global      GlobalConfig      `yaml:"global"`
}

type QBittorrentConfig struct {
	BaseURL            string `yaml:"baseUrl" env:"QBE_URL"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify" env:"QBE_INSECURE_SKIP_VERIFY"`
	Username           string `yaml:"username" env:"QBE_USERNAME"`
	Password           string `yaml:"password" env:"QBE_PASSWORD"`
	Timeout            int    `yaml:"timeout" env:"QBE_TIMEOUT"`
}

type MetricsConfig struct {
	Port    string `yaml:"port" env:"QBE_METRICS_PORT"`
	UrlPath string `yaml:"urlPath" env:"QBE_METRICS_PATH"`
}

type GlobalConfig struct {
	StatePath string `yaml:"statePath" env:"QBE_STATE_PATH"`
}

func UpdatePath(path string) {
	configPath = path
}

func Get() Config {
	mu.Lock()
	defer mu.Unlock()
	if !isLoaded {
		log.Debug("Loading configuration from: " + configPath)
		instance = initializeConfig(configPath)
		isLoaded = true
	}
	return instance
}

func initializeConfig(path string) Config {
	var cfg Config
	if err := validator.ValidatePath(path, false); err != nil {
		panic(err)
	}
	if err := parser.ParseYamlFile(path, &cfg); err != nil {
		panic(fmt.Errorf("error parsing config: %w", err))
	}
	log.Debug("Loading environment variables into configuration")
	loadEnvs(&cfg)
	return cfg
}

func ValidateMetricsPort(cfg Config) error {
	if _, err := strconv.ParseUint(cfg.Metrics.Port, 10, 16); err != nil {
		return fmt.Errorf("invalid metrics port: %w", err)
	}
	return nil
}
