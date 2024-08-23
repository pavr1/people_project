package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func NewConfig(log *log.Logger) (*Config, error) {
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
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		log.Error("AUTH_PORT is not set")
		return nil, errors.New("AUTH_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.WithField("error", err).Error("Failed to convert port to int")
		return nil, err
	}

	var config = Config{}
	config.Server.Port = portInt

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
