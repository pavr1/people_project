package config

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
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
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Error("SERVER_PORT is not set")
		return nil, errors.New("SERVER_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.WithField("error", err).Error("Failed to convert port to int")
		return nil, err
	}

	authPath := os.Getenv("AUTH_PATH")
	if authPath == "" {
		log.Error("AUTH_PATH is not set")
		return nil, errors.New("AUTH_PATH is not set")
	}

	authHost := os.Getenv("AUTH_HOST")
	if authHost == "" {
		log.Error("AUTH_HOST is not set")
		return nil, errors.New("AUTH_HOST is not set")
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		log.Error("MONGODB_URI is not set")
		return nil, errors.New("MONGODB_URI is not set")
	}

	mongodb_database := os.Getenv("MONGODB_DATABASE")
	if mongodb_database == "" {
		log.Error("MONGODB_DATABASE is not set")
		return nil, errors.New("MONGODB_DATABASE is not set")
	}

	mongodb_collection := os.Getenv("MONGODB_COLLECTION")
	if mongodb_collection == "" {
		log.Error("MONGODB_COLLECTION is not set")
		return nil, errors.New("MONGODB_COLLECTION is not set")
	}

	mongodb_username := os.Getenv("MONGODB_USERNAME")
	if mongodb_username == "" {
		log.Error("MONGODB_USERNAME is not set")
		return nil, errors.New("MONGODB_USERNAME is not set")
	}

	mongodb_password := os.Getenv("MONGODB_PASSWORD")
	if mongodb_password == "" {
		log.Error("MONGODB_PASSWORD is not set")
		return nil, errors.New("MONGODB_PASSWORD is not set")
	}

	mongodb_role := os.Getenv("MONGODB_ROLE")
	if mongodb_role == "" {
		log.Error("MONGODB_ROLE is not set")
		return nil, errors.New("MONGODB_ROLE is not set")
	}

	var config = Config{}
	config.Server.Port = portInt
	config.Auth.Path = authPath
	config.Auth.Host = authHost
	config.MongoDB.Uri = mongodb_uri
	config.MongoDB.Database = mongodb_database
	config.MongoDB.Collection = mongodb_collection
	config.MongoDB.Username = mongodb_username
	config.MongoDB.Password = mongodb_password
	config.MongoDB.RolName = mongodb_role

	log.WithField("config", config).Info("Loaded configuration file")

	return &config, nil
}

// Unmarshal the configuration into a struct
// var config Config
// err = viper.Unmarshal(&config)
// if err != nil {
// 	log.WithField("error", err).Error("Failed to unmarshal configuration file")
// 	return nil, err
// }

// log.WithField("config", config).Info("Loaded configuration file")

// if config.Server.Port == 0 {
// 	return nil, errors.New("server.port is required")
// }

// if config.Auth.Path == "" {
// 	return nil, errors.New("auth.path is required")
// }

// if config.Auth.Host == "" {
// 	return nil, errors.New("auth.host is required")
// }

// if config.MongoDB.Uri == "" {
// 	return nil, errors.New("mongodb.uri is required")
// }

// if config.MongoDB.Database == "" {
// 	return nil, errors.New("mongodb.database is required")
// }

// if config.MongoDB.Collection == "" {
// 	return nil, errors.New("mongodb.collection is required")
// }

// if config.MongoDB.Username == "" {
// 	return nil, errors.New("mongodb.username is required")
// }

// if config.MongoDB.Password == "" {
// 	return nil, errors.New("mongodb.password is required")
// }

// if config.MongoDB.RolName == "" {
// 	return nil, errors.New("mongodb.role is required")
// }
