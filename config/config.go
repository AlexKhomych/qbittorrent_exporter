package config

import (
	"fmt"
	"os"
	"qbittorrent_exporter/lib/log"
	"strconv"
)

const (
	e_QBT_USERNAME string = "QBT_USERNAME"
	e_QBT_PASSWORD string = "QBT_PASSWORD"
	e_QBT_BASE_URL string = "QBT_BASE_URL"

	e_METRICS_PORT     string = "METRICS_PORT"
	e_METRICS_URL_PATH string = "METRICS_URL_PATH"
)

type Config struct {
	QBittorrent QBittorrentConfig `yaml:"qBittorrent"`
	Metrics     MetricsConfig     `yaml:"metrics"`
}

type QBittorrentConfig struct {
	BaseURL  string `yaml:"baseUrl"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MetricsConfig struct {
	Port    string `yaml:"port"`
	UrlPath string `yaml:"urlPath"`
}

func (c Config) WithEnvPriority() Config {
	return WithEnvPriority()
}

func WithEnvPriority() Config {
	config := Get()

	config.QBittorrent.BaseURL = getEnv(e_QBT_BASE_URL, config.QBittorrent.BaseURL)
	config.QBittorrent.Username = getEnv(e_QBT_USERNAME, config.QBittorrent.Username)
	config.QBittorrent.Password = getEnv(e_QBT_PASSWORD, config.QBittorrent.Password)

	config.Metrics.Port = getEnv(e_METRICS_PORT, config.Metrics.Port)
	config.Metrics.UrlPath = getEnv(e_METRICS_URL_PATH, config.Metrics.UrlPath)

	update(config)
	return config
}

func Validate() error {
	config := Get()

	if _, err := strconv.ParseUint(config.Metrics.Port, 10, 16); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to validate config.Metrics.Port")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Info("Using environmental variable", "key", key)
		return value
	}
	return fallback
}
