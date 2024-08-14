package config

import (
	"errors"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Prometheus struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"prometheus"`
	Auth struct {
		Path string `mapstructure:"path"`
		Host string `mapstructure:"host"`
	} `mapstructure:"auth"`
	MongoDB struct {
		Uri        string `mapstructure:"uri"`
		Database   string `mapstructure:"database"`
		Collection string `mapstructure:"collection"`
		//pvillalobos add this to a secret later
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		RolName  string `mapstructure:"role"`
	} `mapstructure:"mongodb"`
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

	if config.Server.Port == 0 {
		return nil, errors.New("server.port is required")
	}

	if config.Prometheus.Port == 0 {
		return nil, errors.New("prometheus.port is required")
	}

	if config.Auth.Path == "" {
		return nil, errors.New("auth.path is required")
	}

	if config.Auth.Host == "" {
		return nil, errors.New("auth.host is required")
	}

	if config.MongoDB.Uri == "" {
		return nil, errors.New("mongodb.uri is required")
	}

	if config.MongoDB.Database == "" {
		return nil, errors.New("mongodb.database is required")
	}

	if config.MongoDB.Collection == "" {
		return nil, errors.New("mongodb.collection is required")
	}

	if config.MongoDB.Username == "" {
		return nil, errors.New("mongodb.username is required")
	}

	if config.MongoDB.Password == "" {
		return nil, errors.New("mongodb.password is required")
	}

	if config.MongoDB.RolName == "" {
		return nil, errors.New("mongodb.role is required")
	}

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}
