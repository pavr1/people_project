package config

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Prometheus struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"prometheus"`
}

func NewConfig() (*Config, error) {
	port := os.Getenv("PROMETHEUS_PORT")
	if port == "" {
		log.Error("PROMETHEUS_PORT is not set")
		return nil, errors.New("SERVER_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.WithField("error", err).Error("Failed to convert port to int")
		return nil, err
	}

	var config = Config{}
	config.Prometheus.Port = portInt

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
