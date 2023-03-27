package config

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type ServiceConfig struct {
	APIPort int `envconfig:"PORT" default:"8181"`
}

type DBConfig struct {
	DB_URL string `envconfig:"DB_URL" default:"localhost:27017"`
}

// GetServiceConfig method to fetch the ServiceConfig
func GetServiceConfig(ctx context.Context) ServiceConfig {
	log.Println("Fetching Service configs")
	config := ServiceConfig{}

	err := envconfig.Process("e-commerce", &config)
	if err != nil {
		log.Fatalln(ctx, "Failed fetching service configs")
		panic(err)
	}
	return config
}

// GetDBConfig get db env vars or error
func GetDBConfig(ctx context.Context) DBConfig {
	dbConfig := DBConfig{}
	err := envconfig.Process("e-commerce", &dbConfig)
	if err != nil {
		log.Fatalln(ctx, "Failed fetching db configs")
		panic(err)
	}
	return dbConfig
}
