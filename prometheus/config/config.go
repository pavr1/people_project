package config

import (
	"errors"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Prometheus struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"prometheus"`
}

func NewConfig() (*Config, error) {
	// Set the configuration file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Set the configuration file path
	viper.AddConfigPath(".")

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to read configuration file")
		return nil, err
	}

	// Unmarshal the configuration into a struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.WithField("error", err).Error("Failed to unmarshal configuration file")
		return nil, err
	}

	log.WithField("config", config).Info("Loaded configuration file")

	if config.Prometheus.Port == 0 {
		return nil, errors.New("prometheus.port is required")
	}

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
