package config

type ServiceConfig struct {
	APIPort int `envconfig:"PORT" default:"8181"`
}

type DBConfig struct {
	DB_URL string `envconfig:"DB_URL" default:"localhost:27017"`
}
